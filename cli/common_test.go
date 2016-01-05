package cli

import (
	"testing"

	"github.com/bitrise-io/go-utils/templateutil"
	"github.com/stretchr/testify/require"
)

func Test_createAvailableTemplateFunctions(t *testing.T) {
	plantVars := map[string]string{
		"MyKey1": "my value 1",
	}
	inventory := map[string]string{}

	t.Log("fn: var - First, a successful one, referencing MyKey1")
	evaluatedContent, err := templateutil.EvaluateTemplateStringToString(`{{ var "MyKey1"}}`, inventory,
		createAvailableTemplateFunctions(plantVars))
	require.NoError(t, err)
	require.Equal(t, "my value 1", evaluatedContent)

	t.Log("fn: var - Now reference a Var which doesn't exist - should error")
	evaluatedContent, err = templateutil.EvaluateTemplateStringToString(`{{ var "MyKeyDoesntExist"}}`, inventory,
		createAvailableTemplateFunctions(plantVars))
	require.Error(t, err)
	require.Equal(t, err.Error(), "template: :1:3: executing \"\" at <var \"MyKeyDoesntExis...>: error calling var: No value found for key: MyKeyDoesntExist")
}
