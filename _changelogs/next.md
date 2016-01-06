* new template function: `notEmpty` - template execution fails if the argument string is empty
  * example, checking an Inventory item, making sure it's not empty: `{{ notEmpty .PlantPath }}`
    or in case of a `var`: `{{ var "EmptyKey" | notEmpty }}`; both will terminate the template
    execution if the value is an empty string.
