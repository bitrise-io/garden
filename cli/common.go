package cli

import (
	"fmt"
	"path"

	"text/template"

	"github.com/bitrise-io/go-utils/pathutil"
)

// GardenTemplateInventoryModel ...
type GardenTemplateInventoryModel struct {
	Vars      map[string]string
	TestBool  bool
	PlantID   string
	PlantPath string
}

func createAvailableTemplateFunctions(inventory GardenTemplateInventoryModel) template.FuncMap {
	return template.FuncMap{
		"isOne": func(i int) bool {
			return i == 1
		},
		"var": func(key string) (string, error) {
			val, isFound := inventory.Vars[key]
			if !isFound {
				return "", fmt.Errorf("No value found for key: %s", key)
			}
			return val, nil
		},
		"notEmpty": func(val string) (string, error) {
			if val == "" {
				return "", fmt.Errorf("Value was empty")
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
