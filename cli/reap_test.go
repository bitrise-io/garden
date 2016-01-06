package cli

import (
	"fmt"
	"os"
	"testing"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func Test_reapPlants(t *testing.T) {
	gardenMap, _, err := loadTestGardenMap()
	require.NoError(t, err)

	// clean up the reap-outputs test dir
	err = os.RemoveAll("../_test/reap-outputs/")
	require.NoError(t, err)

	err = reapPlants(gardenMap.FilteredPlantsIDs("", ""), gardenMap,
		ReapCommandParams{
			Command:     "bash",
			CommandArgs: []string{"../_test/reap_test_script.sh"},
		})
	require.NoError(t, err)

	// check reap outputs
	reapOutputFormat := `_GARDEN_PLANT_PATH: %s
_GARDEN_PLANT_ID: %s
MyVar1: %s
MyVar2: %s
IsItAFruit: %s
IsApples: %s
`
	{
		fileContent, err := fileutil.ReadStringFromFile("../_test/reap-outputs/apple-1.txt")
		require.NoError(t, err)
		absPlantPath, err := pathutil.AbsPath("./PLANTROOT/apple-1-dir")
		require.NoError(t, err)
		expectedContent := fmt.Sprintf(reapOutputFormat,
			absPlantPath,
			"apple-1",
			"my value - for var 1",
			"",
			"this is a fruit",
			"yes")
		require.Equal(t, expectedContent, fileContent)
	}
	{
		fileContent, err := fileutil.ReadStringFromFile("../_test/reap-outputs/orange-1.txt")
		require.NoError(t, err)
		absPlantPath, err := pathutil.AbsPath("./PLANTROOT/orange-1-dir")
		require.NoError(t, err)
		expectedContent := fmt.Sprintf(reapOutputFormat,
			absPlantPath,
			"orange-1",
			"my value - for var 1",
			"my value - for var 2",
			"this is a fruit",
			"no")
		require.Equal(t, expectedContent, fileContent)
	}
}
