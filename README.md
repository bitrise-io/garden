# Garden

A tool to manage your template (plan) based directories.

You can perform a setup (plant) by running `garden plant`,
which'll create your garden (directories) based on your plans (temlates).

You can then run scripts on your whole garden: `garden reap bash my_script.sh`
or on a part of it: `garden --part tomatoes reap bash my_script.sh`
or on a single configuration: `garden --id plantid reap bash my_script.sh`.
