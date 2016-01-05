package cli

import (
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"

	"text/template"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/garden/config"
)

func createAvailableTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"isOne": func(i int) bool {
			return i == 1
		},
	}
}

func checkSeedDir(gardenDirAbsPth, seedPath string) (string, error) {
	seedFullPth := path.Join(gardenDirAbsPth, "seeds", seedPath)
	isExist, err := pathutil.IsDirExists(seedFullPth)
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", fmt.Errorf("No Seed directory found at path: %s", seedFullPth)
	}
	return seedFullPth, nil
}

// loadGardenMap ..
func loadGardenMap() (config.GardenMapModel, string, error) {
	relPth, absPth, err := config.FindGardenDirPath()
	if err != nil {
		return config.GardenMapModel{}, "", fmt.Errorf("Failed to find Garden directory: %s", err)
	}
	log.Printf("=> Using Garden directory: %s (abs path: %s)", colorstring.Green(relPth), absPth)

	gardenMapPth := path.Join(absPth, "map.yml")
	gardenMap, err := config.CreateGardenMapModelFromYMLFile(gardenMapPth)
	if err != nil {
		return config.GardenMapModel{}, "", fmt.Errorf("Failed to load Garden Map (path:%s) with error: %s", gardenMapPth, err)
	}
	return gardenMap, absPth, nil
}
