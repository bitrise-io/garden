package cli

import (
	"testing"

	"github.com/bitrise-io/go-utils/templateutil"
	"github.com/stretchr/testify/require"
)

func Test_createAvailableTemplateFunctions(t *testing.T) {
	inventory := GardenTemplateInventoryModel{
		Vars: map[string]string{
			"MyKey1":   "my value 1",
			"EmptyKey": "",
		},
		PlantID: "p1",
	}

	t.Log("fn: var - First, a successful one, referencing MyKey1")
	evaluatedContent, err := templateutil.EvaluateTemplateStringToString(`{{ var "MyKey1"}}`, inventory,
		createAvailableTemplateFunctions(inventory))
	require.NoError(t, err)
	require.Equal(t, "my value 1", evaluatedContent)

	t.Log("fn: var - Now reference a Var which doesn't exist - should error")
	evaluatedContent, err = templateutil.EvaluateTemplateStringToString(`{{ var "MyKeyDoesntExist"}}`, inventory,
		createAvailableTemplateFunctions(inventory))
	require.EqualError(t, err,
		"template: :1:3: executing \"\" at <var \"MyKeyDoesntExis...>: error calling var: No value found for key: MyKeyDoesntExist")

	// fn: isOne
	t.Log("fn: isOne - should be true")
	evaluatedContent, err = templateutil.EvaluateTemplateStringToString(`{{ isOne 1 }}`, inventory,
		createAvailableTemplateFunctions(inventory))
	require.NoError(t, err)
	require.Equal(t, "true", evaluatedContent)

	t.Log("fn: isOne - should be false")
	evaluatedContent, err = templateutil.EvaluateTemplateStringToString(`{{ isOne 2 }}`, inventory,
		createAvailableTemplateFunctions(inventory))
	require.NoError(t, err)
	require.Equal(t, "false", evaluatedContent)

	// fn: notEmpty
	t.Log("fn: notEmpty - should be ok")
	evaluatedContent, err = templateutil.EvaluateTemplateStringToString(`{{ notEmpty .PlantID }}`, inventory,
		createAvailableTemplateFunctions(inventory))
	require.NoError(t, err)
	require.Equal(t, "p1", evaluatedContent)

	t.Log("fn: notEmpty - should return an error, in case the value is empty")
	evaluatedContent, err = templateutil.EvaluateTemplateStringToString(`{{ notEmpty .PlantPath }}`, inventory,
		createAvailableTemplateFunctions(inventory))
	require.EqualError(t, err,
		"template: :1:3: executing \"\" at <notEmpty .PlantPath>: error calling notEmpty: Value was empty")

	t.Log("fn: notEmpty - combined with var - should return an error, in case the value is empty")
	evaluatedContent, err = templateutil.EvaluateTemplateStringToString(`{{ var "EmptyKey" | notEmpty }}`, inventory,
		createAvailableTemplateFunctions(inventory))
	require.EqualError(t, err,
		"template: :1:20: executing \"\" at <notEmpty>: error calling notEmpty: Value was empty")
}
