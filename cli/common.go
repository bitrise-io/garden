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

func createAvailableTemplateFunctions(plantVars map[string]string) template.FuncMap {
	return template.FuncMap{
		"isOne": func(i int) bool {
			return i == 1
		},
		"var": func(key string) (string, error) {
			val, isFound := plantVars[key]
			if !isFound {
				return "", fmt.Errorf("No value found for key: %s", key)
			}
			return val, nil
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
//  gardenDirPath is optional, if provided will be used as the Garden Dir path
//  if not provided the standard .garden dir paths will be checked
func loadGardenMap(gardenDirPath string) (config.GardenMapModel, string, error) {
	relPath := ""
	absPath := ""

	if gardenDirPath != "" {
		relPath = gardenDirPath
		apth, err := pathutil.AbsPath(gardenDirPath)
		if err != nil {
			return config.GardenMapModel{}, "", fmt.Errorf("Failed to get Absolute path of provided Garden Dir (path:%s), error: %s", gardenDirPath, err)
		}
		absPath = apth
	} else {
		rpth, apth, err := config.FindGardenDirPath()
		if err != nil {
			return config.GardenMapModel{}, "", fmt.Errorf("Failed to find Garden directory: %s", err)
		}
		relPath = rpth
		absPath = apth
	}
	log.Printf("=> Using Garden directory: %s (abs path: %s)", colorstring.Green(relPath), absPath)

	gardenMapPth := path.Join(absPath, "map.yml")
	gardenMap, err := config.CreateGardenMapModelFromYMLFile(gardenMapPth)
	if err != nil {
		return config.GardenMapModel{}, "", fmt.Errorf("Failed to load Garden Map (path:%s) with error: %s", gardenMapPth, err)
	}
	return gardenMap, absPath, nil
}
