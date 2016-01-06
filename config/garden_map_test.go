package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testGardenDirPath = "../_test/garden"
)

func loadTestGardenMap() (GardenMapModel, string, error) {
	gardenMap, gardenDirAbsPth, err := LoadGardenMap(testGardenDirPath)
	if err != nil {
		return GardenMapModel{}, "", fmt.Errorf("Failed to load Garden Map: %s", err)
	}

	return gardenMap, gardenDirAbsPth, nil
}

func Test_GardenMapModel_GetAllVarsForPlant(t *testing.T) {
	gardenMap, _, err := loadTestGardenMap()
	require.NoError(t, err)

	allVars, err := gardenMap.CollectAllVarsForPlant("orange-1")
	require.NoError(t, err)
	require.EqualValues(t,
		map[string]string{
			"MyVar1":     "my value - for var 1",
			"MyVar2":     "my value - for var 2",
			"IsItAFruit": "this is a fruit",
			"IsApples":   "no",
		},
		allVars)

	allVars, err = gardenMap.CollectAllVarsForPlant("apple-1")
	require.NoError(t, err)
	require.EqualValues(t,
		map[string]string{
			"MyVar1":     "my value - for var 1",
			"IsItAFruit": "this is a fruit",
			"IsApples":   "yes",
		},
		allVars)
}

func Test_PlantModel_ExpandedPath(t *testing.T) {
	t.Log("Simple Path")
	plant := PlantModel{
		Path: "abc/def",
	}
	require.Equal(t, "abc/def", plant.ExpandedPath("PLANT1"))

	t.Log("Single expand")
	plant = PlantModel{
		Path: "abc/$_GARDEN_PLANT_ID",
	}
	require.Equal(t, "abc/PLANT1", plant.ExpandedPath("PLANT1"))

	t.Log("Multiple expands")
	plant = PlantModel{
		Path: "abc/$_GARDEN_PLANT_ID/a/$_GARDEN_PLANT_ID",
	}
	require.Equal(t, "abc/PLANT1/a/PLANT1", plant.ExpandedPath("PLANT1"))
}
