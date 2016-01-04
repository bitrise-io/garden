package config

import (
	"errors"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
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

// FilteredPlants ...
func (gardenMap GardenMapModel) FilteredPlants(plantID, zone string) map[string]PlantModel {
	if plantID != "" {
		plantModel, isFound := gardenMap.Plants[plantID]
		if isFound {
			return map[string]PlantModel{
				plantID: plantModel,
			}
		}
		// not found by ID
		return map[string]PlantModel{}
	}

	if zone != "" {
		return gardenMap.plantsFilteredByZone(zone)
	}

	// no filter
	return gardenMap.Plants
}

func (gardenMap GardenMapModel) plantsFilteredByZone(zone string) map[string]PlantModel {
	filtered := map[string]PlantModel{}
	for plantID, plantModel := range gardenMap.Plants {
		if sliceutil.IndexOfStringInSlice(zone, plantModel.Zones) >= 0 {
			filtered[plantID] = plantModel
		}
	}
	return filtered
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
