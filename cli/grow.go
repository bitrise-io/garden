package cli

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/garden/config"
	"github.com/codegangsta/cli"
)

func growPlant(gardenDirAbsPth, plantID string, plantModel config.PlantModel) error {
	log.Println(colorstring.Yellow("=>"), "growing plant:", colorstring.Green(plantID))

	log.Println("--> Checking seed: ", plantModel.Seed, "...")
	seedDirFullPth, err := checkSeedDir(gardenDirAbsPth, plantModel.Seed)
	if err != nil {
		return fmt.Errorf("Failed to check seed directory: %s", err)
	}
	tmpSeedPth, err := pathutil.NormalizedOSTempDirPath("")
	log.Debugln(" tmpSeedPth: ", tmpSeedPth)
	if err != nil {
		return fmt.Errorf("Failed to create a temporary directory for seed: %s", err)
	}
	output, err := cmdex.RunCommandAndReturnCombinedStdoutAndStderr("rsync",
		"-avhP", seedDirFullPth, tmpSeedPth)
	if err != nil {
		log.Errorf("Failed to rsync seed to temporary seed dir: %s", err)
		log.Errorf("Output was: %s", output)
		return err
	}

	log.Println("--> Handling templates ...")

	log.Println("--> Moving plant to it's final place in the garden ...")

	log.Println("--> Cleaning up ...")
	if err := os.RemoveAll(tmpSeedPth); err != nil {
		return fmt.Errorf("Failed to cleanup: %s", err)
	}

	log.Println("-> Plant grown! 🌴")
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

	if err := growPlants(gardenDirAbsPth, gardenMap.FilteredPlants(WorkWithPlantID, WorkWithZone)); err != nil {
		log.Fatalf("Failed to grow plants: %s", err)
	}
}
