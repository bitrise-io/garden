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

func reapThisPlant(plantID string, gardenMap config.GardenMapModel, cmdParams ReapCommandParams) error {
	plant, isFound := gardenMap.Plants[plantID]
	if !isFound {
		return fmt.Errorf("reapThisPlant: can't find Plant with ID: %s", plantID)
	}

	expandedPlantPath := plant.ExpandedPath(plantID)
	absPlantDirPath, err := pathutil.AbsPath(expandedPlantPath)
	if err != nil {
		return fmt.Errorf("Failed to get Absolute Path of Plant (path:%s), error: %s", expandedPlantPath, err)
	}

	cmd := exec.Command(cmdParams.Command, cmdParams.CommandArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// envs
	envsToAdd := []string{}
	// defaults
	envsToAdd = append(envsToAdd, fmt.Sprintf("_GARDEN_PLANT_DIR=%s", absPlantDirPath))
	envsToAdd = append(envsToAdd, fmt.Sprintf("_GARDEN_PLANT_ID=%s", plantID))
	// Vars
	allPlantVars, err := gardenMap.CollectAllVarsForPlant(plantID)
	if err != nil {
		return fmt.Errorf("reapThisPlant: failed to collect Vars for Plant (id: %s), error: %s", plantID, err)
	}
	log.Debugf("allPlantVars: %#v", allPlantVars)
	for key, val := range allPlantVars {
		envsToAdd = append(envsToAdd, fmt.Sprintf("_GARDENVAR_%s=%s", key, val))
	}
	cmd.Env = append(os.Environ(), envsToAdd...)

	return cmd.Run()
}

func reapPlants(plantIDs []string, gardenMap config.GardenMapModel, cmdParams ReapCommandParams) error {
	for _, plantID := range plantIDs {
		fmt.Println()
		log.Infof("🚜  -> Reaping plant: %s", colorstring.Green(plantID))
		if err := reapThisPlant(plantID, gardenMap, cmdParams); err != nil {
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

	plantsToGrowIDs := gardenMap.FilteredPlantsIDs(WorkWithPlantID, WorkWithZone)
	if len(plantsToGrowIDs) < 1 {
		log.Fatalln("No plants to grow!")
	}
	if err := reapPlants(plantsToGrowIDs, gardenMap, cmdParams); err != nil {
		log.Fatalf("Failed to grow plants: %s", err)
	}
}
