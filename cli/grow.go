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
	"github.com/bitrise-io/garden/config"
	"github.com/codegangsta/cli"
)

// evaluateAndReplaceTemplateFile ...
//  it'll evalutate the content of the template file
//  and then write it into a new file, without the .template extension
//  and remove the original template file
func evaluateAndReplaceTemplateFile(templateFilePath string, templateInventory GardenTemplateInventoryModel) error {
	fileContent, err := fileutil.ReadStringFromFile(templateFilePath)
	if err != nil {
		return fmt.Errorf("Failed to read template file (path:%s), error: %s", templateFilePath, err)
	}
	evaluatedContent, err := templateutil.EvaluateTemplateStringToString(fileContent, templateInventory,
		createAvailableTemplateFunctions(templateInventory))
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

func replaceTemplateFilesInDir(dirPth string, templateInventory GardenTemplateInventoryModel) error {
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

	for _, aTemplateFilePth := range templateFilePaths {
		log.Infoln(colorstring.Cyan("-> Evaluating and replacing template file:"), aTemplateFilePth)
		if err := evaluateAndReplaceTemplateFile(aTemplateFilePth, templateInventory); err != nil {
			return fmt.Errorf("Failed to evaluate template file (path:%s), error: %s", aTemplateFilePth, err)
		}
	}

	return nil
}

func growPlant(plantID string, gardenMap config.GardenMapModel, gardenDirAbsPth string) error {
	fmt.Println()
	log.Println(colorstring.Yellow("==> growing plant:"), colorstring.Green(plantID))
	log.Println("🌱")

	plantModel, isFound := gardenMap.Plants[plantID]
	if !isFound {
		return fmt.Errorf("growPlant: can't find Plant with ID: %s", plantID)
	}

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
	expandedPlantPath := plantModel.ExpandedPath(plantID)
	absPlantPath, err := pathutil.AbsPath(expandedPlantPath)
	if err != nil {
		return fmt.Errorf("Failed to get Absolute path of plant (path:%s), error: %s", expandedPlantPath, err)
	}
	collectedPlantVars, err := gardenMap.CollectAllVarsForPlant(plantID)
	if err != nil {
		return fmt.Errorf("growPlant: failed to collect Vars for Plant (id: %s), error: %s", plantID, err)
	}
	templateInventory := GardenTemplateInventoryModel{
		TestBool:  true,
		Vars:      collectedPlantVars,
		PlantID:   plantID,
		PlantPath: absPlantPath,
	}

	if err := replaceTemplateFilesInDir(tmpSeedPth, templateInventory); err != nil {
		return fmt.Errorf("Failed to handle templates in temp seed dir (path:%s), error: %s", tmpSeedPth, err)
	}

	log.Println("--> Moving plant to it's final place in the garden ...")
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

func growPlants(gardenDirAbsPth string, gardenMap config.GardenMapModel, plantsToGrowIDs []string) error {
	for _, plantID := range plantsToGrowIDs {
		if err := growPlant(plantID, gardenMap, gardenDirAbsPth); err != nil {
			return err
		}
	}
	return nil
}

func grow(c *cli.Context) {
	log.Infoln("Grow")

	gardenMap, gardenDirAbsPth, err := config.LoadGardenMap("")
	if err != nil {
		log.Fatalf("Failed to load Garden Map: %s", err)
	}

	plantsToGrowIDs := gardenMap.FilteredPlantsIDs(WorkWithPlantID, WorkWithZone)
	if len(plantsToGrowIDs) < 1 {
		log.Fatalln("No plants to grow!")
	}
	if err := growPlants(gardenDirAbsPth, gardenMap, plantsToGrowIDs); err != nil {
		log.Fatalf("Failed to grow plants: %s", err)
	}
}
