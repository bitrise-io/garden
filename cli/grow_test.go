package cli

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/garden/config"
	"github.com/stretchr/testify/require"
)

const (
	testGardenDirPath          = "../_test/garden"
	testGardenPlantrootRelPath = "../_test/plantroot"
)

func loadTestGardenMap() (config.GardenMapModel, string, error) {
	gardenMap, gardenDirAbsPth, err := config.LoadGardenMap(testGardenDirPath)
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

func testFileContent(t *testing.T, filePth, expectedContent string) {
	filecont, err := fileutil.ReadStringFromFile(filePth)
	require.NoError(t, err)
	require.Equal(t, expectedContent, filecont)
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
	t.Logf("-> gardenMap: %#v", gardenMap)

	err = growPlants(absTestGardenDirPath, gardenMap, gardenMap.FilteredPlantsIDs("", ""))
	require.NoError(t, err)

	// Apple-1
	// test the generated files
	appleOneDirPth := path.Join(absPlantRootPath, "apple-1-dir")
	// fix file, no template content
	testFileContent(t, path.Join(appleOneDirPth, "fix-file.txt"), `Apples - this is a non template file.
`)
	// template 1, in root dir of plant
	testFileContent(t, path.Join(appleOneDirPth, "templated-file.txt"), `Apples - this is a templated file.

Temp:  T1 |
Value of MyVar1: my value - for var 1
IsItAFruit: this is a fruit
IsApples: yes
`)
	// template 2, in a subdir of plant
	testFileContent(t, path.Join(appleOneDirPth, "subdir", "tempinsub"), `Apples - this is a templated file, in a sub directory.

Temp:  T1 |
Temp: |
Value of MyVar1: my value - for var 1
`)

	// Orange-1
	// test the generated files
	orangeOneDirPth := path.Join(absPlantRootPath, "orange-1-dir")
	// fix file, no template content
	testFileContent(t, path.Join(orangeOneDirPth, "fix-file.txt"), `Oranges - this is a non template file.
`)
	// template 1, in root dir of plant
	testFileContent(t, path.Join(orangeOneDirPth, "templated-file.txt"), `Oranges - this is a templated file.

Temp:  T1 |
Value of MyVar1: my value - for var 1
Value of MyVar2: my value - for var 2
IsItAFruit: this is a fruit
IsApples: no
`)

}
