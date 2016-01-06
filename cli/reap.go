package cli

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/garden/config"
	"github.com/codegangsta/cli"
)

// ReapCommandParams ...
type ReapCommandParams struct {
	Command     string
	CommandArgs []string
}

func reapThisPlant(plant config.PlantModel, cmdParams ReapCommandParams) error {
	absPlantDirPath, err := pathutil.AbsPath(plant.Path)
	if err != nil {
		return fmt.Errorf("Failed to get Absolute Path of Plant (path:%s), error: %s", plant.Path, err)
	}

	cmd := exec.Command(cmdParams.Command, cmdParams.CommandArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// envs
	envsToAdd := []string{}
	// defaults
	envsToAdd = append(envsToAdd, fmt.Sprintf("_GARDEN_PLANTDIR=%s", absPlantDirPath))
	// Vars
	for key, val := range plant.Vars {
		envsToAdd = append(envsToAdd, fmt.Sprintf("_GARDENVAR_%s=%s", key, val))
	}
	cmd.Env = append(os.Environ(), envsToAdd...)

	return cmd.Run()
}

func reapPlants(plants config.PlantsMap, cmdParams ReapCommandParams) error {
	for plantID, plantModel := range plants {
		fmt.Println()
		log.Infof("🚜  -> Reaping plant: %s", colorstring.Green(plantID))
		if err := reapThisPlant(plantModel, cmdParams); err != nil {
			return fmt.Errorf("Failed to reap plant (id:%s), error: %s", plantID, err)
		}
	}
	return nil
}

func reap(c *cli.Context) {
	log.Infoln("Reap")

	args := c.Args()
	log.Infof("args: %#v", args)
	if len(args) < 1 {
		log.Fatalln("No command to execute, can't reap.")
	}
	cmdParams := ReapCommandParams{
		Command: args[0],
	}
	if len(args) > 1 {
		cmdParams.CommandArgs = args[1:]
	}

	gardenMap, _, err := config.LoadGardenMap("")
	if err != nil {
		log.Fatalf("Failed to load Garden Map: %s", err)
	}

	plantsToGrow := gardenMap.FilteredPlants(WorkWithPlantID, WorkWithZone)
	if len(plantsToGrow) < 1 {
		log.Fatalln("No plants to grow!")
	}
	if err := reapPlants(plantsToGrow, cmdParams); err != nil {
		log.Fatalf("Failed to grow plants: %s", err)
	}
}
