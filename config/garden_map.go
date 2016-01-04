package config

import (
	"errors"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"gopkg.in/yaml.v2"
)

// PlantModel ...
type PlantModel struct {
	Path  string   `json:"path" yaml:"path"`
	Seed  string   `json:"seed" yaml:"seed"`
	Vars  []string `json:"vars" yaml:"vars"`
	Zones []string `json:"zones" yaml:"zones"`
}

// GardenMapModel ...
type GardenMapModel struct {
	Plants map[string]PlantModel `json:"plants" yaml:"plants"`
}

func checkGardenDirPath(relPth string) (string, string, error) {
	isEx, err := pathutil.IsDirExists(relPth)
	if err != nil {
		return "", "", err
	}
	if isEx {
		absPth, err := pathutil.AbsPath(relPth)
		if err != nil {
			return "", "", err
		}
		return relPth, absPth, nil
	}
	return "", "", nil
}

// FindGardenDirPath ...
func FindGardenDirPath() (string, string, error) {
	relPth := "./.garden"
	if _, absPth, err := checkGardenDirPath(relPth); err == nil && absPth != "" {
		return relPth, absPth, nil
	}
	relPth = "~/.garden"
	if _, absPth, err := checkGardenDirPath(relPth); err == nil && absPth != "" {
		return relPth, absPth, nil
	}
	return "", "", errors.New("Can't find Garden directory at standard paths")
}

// CreateGardenMapModelFromYMLFile ...
func CreateGardenMapModelFromYMLFile(pth string) (GardenMapModel, error) {
	fileBytes, err := fileutil.ReadBytesFromFile(pth)
	if err != nil {
		return GardenMapModel{}, err
	}

	var modelToReturn GardenMapModel
	if err := yaml.Unmarshal(fileBytes, &modelToReturn); err != nil {
		return GardenMapModel{}, err
	}

	return modelToReturn, nil
}
