* __FIX__ : `garden` will now find and be able to use the Garden Dir
  if it's located at `~/.garden`, not just in case the `.garden` dir is
  located in the current directory
* `PlantPath` is now available in the template Inventory
  * Usage example: `{{ .PlantPath }}`
