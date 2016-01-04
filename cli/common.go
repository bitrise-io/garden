package cli

import (
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-tools/garden/config"
)

// loadGardenMap ..
func loadGardenMap() (config.GardenMapModel, error) {
	relPth, absPth, err := config.FindGardenDirPath()
	if err != nil {
		return config.GardenMapModel{}, fmt.Errorf("Failed to find Garden directory: %s", err)
	}
	log.Printf("=> Using Garden directory: %s (abs path: %s)", colorstring.Green(relPth), absPth)

	gardenMapPth := path.Join(absPth, "map.yml")
	gardenMap, err := config.CreateGardenMapModelFromYMLFile(gardenMapPth)
	if err != nil {
		return config.GardenMapModel{}, fmt.Errorf("Failed to load Garden Map (path:%s) with error: %s", gardenMapPth, err)
	}
	return gardenMap, nil
}
