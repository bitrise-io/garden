package cli

import (
	"fmt"
	"path"

	"text/template"

	"github.com/bitrise-io/go-utils/pathutil"
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
