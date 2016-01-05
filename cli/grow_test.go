package cli

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/garden/config"
	"github.com/stretchr/testify/require"
)

const (
	testGardenDirPath          = "../_test/garden"
	testGardenPlantrootRelPath = "../_test/plantroot"
)

func loadTestGardenMap() (config.GardenMapModel, string, error) {
	gardenMap, gardenDirAbsPth, err := loadGardenMap(testGardenDirPath)
	if err != nil {
		return config.GardenMapModel{}, "", fmt.Errorf("Failed to load Garden Map: %s", err)
	}

	return gardenMap, gardenDirAbsPth, nil
}

func fixPlantPathForTest(gardenMap config.GardenMapModel, plantrootPath string) config.GardenMapModel {
	for key, val := range gardenMap.Plants {
		val.Path = strings.Replace(val.Path, "PLANTROOT", plantrootPath, 1)
		gardenMap.Plants[key] = val
	}
	return gardenMap
}

func Test_growPlants(t *testing.T) {
	gardenMap, absTestGardenDirPath, err := loadTestGardenMap()
	require.NoError(t, err)

	// replace PLANTROOT to point to our _test/plantroot
	absPlantRootPath, err := pathutil.AbsPath(testGardenPlantrootRelPath)
	require.NoError(t, err)
	{
		isEx, err := pathutil.IsDirExists(absPlantRootPath)
		require.NoError(t, err)
		if isEx {
			err = os.RemoveAll(absPlantRootPath)
			require.NoError(t, err)
		}
	}
	err = os.Mkdir(absPlantRootPath, 0777)
	require.NoError(t, err)

	gardenMap = fixPlantPathForTest(gardenMap, absPlantRootPath)
	log.Printf("-> gardenMap: %#v", gardenMap)

	err = growPlants(absTestGardenDirPath, gardenMap.Plants)
	require.NoError(t, err)
}
