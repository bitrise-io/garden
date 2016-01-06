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
