# Garden

A tool to manage your template (plan) based directories.

You can perform a setup (plant) by running `garden grow`,
which'll create your garden (directories) based on your plans (temlates).

You can filter every command with `-zone=` and `-plant`, so to grow
only a part of your garden you can: `garden -zone=tomatoes grow`
or to grow a single plant: `garden -plant=fav-tomatoe grow`

You can then run scripts on your whole garden: `garden reap bash my_script.sh`
or on a part of it: `garden -zone=tomatoes reap bash my_script.sh`
or on a single configuration: `garden -plant=fav-tomatoe reap bash my_script.sh`.

Reap can access the Vars of the Plant as Environment Variables,
in the following form: `_GARDENVAR_[the-Var-id]`, as well as
a couple of built-in Environment Variables:

* `_GARDEN_PLANT_DIR` : the Absolute Directory Path of the Plant
* `_GARDEN_PLANT_ID` : ID of the Plant

You can test & view your garden with `garden view`.
