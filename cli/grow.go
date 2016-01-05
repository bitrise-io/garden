package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/templateutil"
	"github.com/bitrise-tools/garden/config"
	"github.com/codegangsta/cli"
)

// GrowInventoryModel ...
type GrowInventoryModel struct {
	Vars     map[string]string
	TestBool bool
}

// evaluateAndReplaceTemplateFile ...
//  it'll evalutate the content of the template file
//  and then write it into a new file, without the .template extension
//  and remove the original template file
func evaluateAndReplaceTemplateFile(templateFilePath string, templateInventory GrowInventoryModel) error {
	fileContent, err := fileutil.ReadStringFromFile(templateFilePath)
	if err != nil {
		return fmt.Errorf("Failed to read template file (path:%s), error: %s", templateFilePath, err)
	}
	evaluatedContent, err := templateutil.EvaluateTemplateStringToString(fileContent, templateInventory, createAvailableTemplateFunctions())
	if err != nil {
		return fmt.Errorf("Failed to evaluate template (path:%s), error: %s", templateFilePath, err)
	}

	evaluatedFileSavePth := strings.TrimSuffix(templateFilePath, ".template")
	origFilePerms, err := fileutil.GetFilePermissions(templateFilePath)
	if err != nil {
		return fmt.Errorf("Failed to get permission settings of the original template file, error: %s", err)
	}
	log.Println("Writing evaluated template content into file:", evaluatedFileSavePth)

	if err := fileutil.WriteStringToFileWithPermission(evaluatedFileSavePth, evaluatedContent, origFilePerms); err != nil {
		return fmt.Errorf("Failed to write evaluated content into file (path:%s), error: %s", evaluatedFileSavePth, err)
	}

	if err := os.Remove(templateFilePath); err != nil {
		return fmt.Errorf("Failed to delete the temporary template file (path:%s), error: %s", templateFilePath, err)
	}

	return nil
}

func replaceTemplateFilesInDir(dirPth string) error {
	templateFilePaths := []string{}
	err := filepath.Walk(dirPth, func(pth string, f os.FileInfo, err error) error {
		if f.Mode().IsDir() {
			log.Debugf("-> (i) Path is directory, skipping: %s", pth)
			return nil
		}
		log.Debugf("-> Checking path: %s / ext: %s", pth, filepath.Ext(pth))
		if filepath.Ext(pth) == ".template" {
			log.Debugln(colorstring.Cyanf("--> Template Found! : %s", pth))
			templateFilePaths = append(templateFilePaths, pth)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to scan template files in directory (path:%s), error: %s", dirPth, err)
	}

	log.Infoln(colorstring.Cyan("-> templateFilePaths:"), templateFilePaths)

	templateInventory := GrowInventoryModel{TestBool: true}

	for _, aTemplateFilePth := range templateFilePaths {
		log.Infoln(colorstring.Cyan("-> Evaluating and replacing template file:"), aTemplateFilePth)
		if err := evaluateAndReplaceTemplateFile(aTemplateFilePth, templateInventory); err != nil {
			return fmt.Errorf("Failed to evaluate template file (path:%s), error: %s", aTemplateFilePth, err)
		}
	}

	return nil
}

func growPlant(gardenDirAbsPth, plantID string, plantModel config.PlantModel) error {
	fmt.Println()
	log.Println(colorstring.Yellow("==> growing plant:"), colorstring.Green(plantID))
	log.Println("🌱")

	log.Println("--> Checking seed: ", plantModel.Seed, "...")
	seedDirFullPth, err := checkSeedDir(gardenDirAbsPth, plantModel.Seed)
	if err != nil {
		return fmt.Errorf("Failed to check seed directory: %s", err)
	}
	tmpSeedPth, err := pathutil.NormalizedOSTempDirPath("")
	log.Debugln("    temp seed dir: ", tmpSeedPth)
	if err != nil {
		return fmt.Errorf("Failed to create a temporary directory for seed: %s", err)
	}
	// only content of dir
	output, err := cmdex.RunCommandAndReturnCombinedStdoutAndStderr("rsync",
		"-avhP", filepath.Clean(seedDirFullPth)+"/", filepath.Clean(tmpSeedPth)+"/")
	if err != nil {
		log.Errorf("Failed to rsync seed to temporary seed dir: %s", err)
		log.Errorf("Output was: %s", output)
		return err
	}

	log.Println("--> Handling templates ...")
	if err := replaceTemplateFilesInDir(tmpSeedPth); err != nil {
		return fmt.Errorf("Failed to handle templates in temp seed dir (path:%s), error: %s", tmpSeedPth, err)
	}

	log.Println("--> Moving plant to it's final place in the garden ...")
	absPlantPath, err := pathutil.AbsPath(plantModel.Path)
	if err != nil {
		return fmt.Errorf("Failed to get Absolute path of plant (path:%s), error: %s", plantModel.Path, err)
	}
	log.Println("    Plant's final place: ", absPlantPath)
	// only content of dir
	output, err = cmdex.RunCommandAndReturnCombinedStdoutAndStderr("rsync",
		"-avhP", filepath.Clean(tmpSeedPth)+"/", filepath.Clean(absPlantPath)+"/")
	if err != nil {
		log.Errorf("Failed to rsync temporary seed dir to it's final place: %s", err)
		log.Errorf("Output was: %s", output)
		return err
	}

	log.Println("--> Cleaning up ...")
	if err := os.RemoveAll(tmpSeedPth); err != nil {
		return fmt.Errorf("Failed to cleanup: %s", err)
	}
	log.Debugln("    [OK] Removed temp seed dir:", tmpSeedPth)

	log.Println("🌴")
	log.Println("-> Plant grown!")
	return nil
}

func growPlants(gardenDirAbsPth string, plantsMap config.PlantsMap) error {
	for plantID, plantModel := range plantsMap {
		if err := growPlant(gardenDirAbsPth, plantID, plantModel); err != nil {
			return err
		}
	}
	return nil
}

func grow(c *cli.Context) {
	log.Infoln("Grow")

	gardenMap, gardenDirAbsPth, err := loadGardenMap()
	if err != nil {
		log.Fatalf("Failed to load Garden Map: %s", err)
	}

	plantsToGrow := gardenMap.FilteredPlants(WorkWithPlantID, WorkWithZone)
	if len(plantsToGrow) < 1 {
		log.Fatalln("No plants to grow!")
	}
	if err := growPlants(gardenDirAbsPth, plantsToGrow); err != nil {
		log.Fatalf("Failed to grow plants: %s", err)
	}
}
