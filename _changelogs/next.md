* `PlantID` is now available in the template Inventory
  * this means that in addition to using Vars (`{{ var "MyVar" }}` or (`{{ .Vars.MyVar }}`)
    you can now reference the PlantID as well: `{{ .PlantID }}`
* Additionally the PlantID can be referenced in the `path` of a plant as well
  with `$_GARDEN_PLANT_ID`
  * Example: `path: ~/my_plants/$_GARDEN_PLANT_ID`
