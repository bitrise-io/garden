package config

import (
	"errors"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/sliceutil"
	"gopkg.in/yaml.v2"
)

// PlantVarsMap ...
type PlantVarsMap map[string]string

// PlantModel ...
type PlantModel struct {
	Path  string       `json:"path" yaml:"path"`
	Seed  string       `json:"seed" yaml:"seed"`
	Vars  PlantVarsMap `json:"vars" yaml:"vars"`
	Zones []string     `json:"zones" yaml:"zones"`
}

// ZoneModel ..
type ZoneModel struct {
	Vars PlantVarsMap `json:"vars" yaml:"vars"`
}

// PlantsMap ...
type PlantsMap map[string]PlantModel

// GardenMapModel ...
type GardenMapModel struct {
	Plants map[string]PlantModel `json:"plants" yaml:"plants"`
	Zones  map[string]ZoneModel  `json:"zones" yaml:"zones"`
}

// ExpandedPath ...
func (plant PlantModel) ExpandedPath(plantID string) string {
	return strings.Replace(plant.Path, "$_GARDEN_PLANT_ID", plantID, -1)
}

// CollectAllVarsForPlant ...
//  collects all the Vars for a plant, including the ones defined for
//  the plant's zones
// In case a variable is defined in multiple Zones or in a Zone and in the Plant
//  as well: Plant's Vars will always be the #1 priority, no matter whether
//  it's also defined in a Zone or not.
//  In case of different Zones, the last one in the Plant's Zones list
//  will be the one which's Var will be used, it'll overwrite other
//  zones' previously defined Vars for the same key.
func (gardenMap GardenMapModel) CollectAllVarsForPlant(plantID string) (PlantVarsMap, error) {
	allVars := PlantVarsMap{}
	plantModel, isFound := gardenMap.Plants[plantID]
	if !isFound {
		return PlantVarsMap{}, fmt.Errorf("Failed to find Plant with ID: %s", plantID)
	}

	for _, aZoneID := range plantModel.Zones {
		zoneModel, isFound := gardenMap.Zones[aZoneID]
		if !isFound {
			// no Zone specific data/vars
			continue
		}
		for k, v := range zoneModel.Vars {
			allVars[k] = v
		}
	}

	for k, v := range plantModel.Vars {
		allVars[k] = v
	}

	return allVars, nil
}

// FilteredPlants ...
func (gardenMap GardenMapModel) FilteredPlants(plantID, zone string) PlantsMap {
	if plantID != "" {
		plantModel, isFound := gardenMap.Plants[plantID]
		if isFound {
			return PlantsMap{
				plantID: plantModel,
			}
		}
		// not found by ID
		return PlantsMap{}
	}

	if zone != "" {
		return gardenMap.plantsFilteredByZone(zone)
	}

	// no filter
	return gardenMap.Plants
}

// FilteredPlantsIDs ...
func (gardenMap GardenMapModel) FilteredPlantsIDs(plantID, zone string) []string {
	filteredPlants := gardenMap.FilteredPlants(plantID, zone)
	ids := []string{}
	for aPlantID := range filteredPlants {
		ids = append(ids, aPlantID)
	}
	return ids
}

func (gardenMap GardenMapModel) plantsFilteredByZone(zone string) PlantsMap {
	filtered := PlantsMap{}
	for plantID, plantModel := range gardenMap.Plants {
		if sliceutil.IndexOfStringInSlice(zone, plantModel.Zones) >= 0 {
			filtered[plantID] = plantModel
		}
	}
	return filtered
}

func checkGardenDirPath(relPth string) (string, string, error) {
	absPth, err := pathutil.AbsPath(relPth)
	if err != nil {
		return "", "", err
	}

	isEx, err := pathutil.IsDirExists(absPth)
	if err != nil {
		return "", "", err
	}
	if !isEx {
		return "", "", nil
	}
	return relPth, absPth, nil
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

// LoadGardenMap ..
//  gardenDirPath is optional, if provided will be used as the Garden Dir path
//  if not provided the standard .garden dir paths will be checked
func LoadGardenMap(gardenDirPath string) (GardenMapModel, string, error) {
	relPath := ""
	absPath := ""

	if gardenDirPath != "" {
		relPath = gardenDirPath
		apth, err := pathutil.AbsPath(gardenDirPath)
		if err != nil {
			return GardenMapModel{}, "", fmt.Errorf("Failed to get Absolute path of provided Garden Dir (path:%s), error: %s", gardenDirPath, err)
		}
		absPath = apth
	} else {
		rpth, apth, err := FindGardenDirPath()
		if err != nil {
			return GardenMapModel{}, "", fmt.Errorf("Failed to find Garden directory: %s", err)
		}
		relPath = rpth
		absPath = apth
	}
	log.Printf("=> Using Garden directory: %s (abs path: %s)", colorstring.Green(relPath), absPath)

	gardenMapPth := path.Join(absPath, "map.yml")
	gardenMap, err := CreateGardenMapModelFromYMLFile(gardenMapPth)
	if err != nil {
		return GardenMapModel{}, "", fmt.Errorf("Failed to load Garden Map (path:%s) with error: %s", gardenMapPth, err)
	}
	return gardenMap, absPath, nil
}
