#!/bin/bash

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

mkdir "$THIS_SCRIPT_DIR/reap-outputs"

(
  echo "_GARDEN_PLANT_PATH: $_GARDEN_PLANT_PATH"
  echo "_GARDEN_PLANT_ID: $_GARDEN_PLANT_ID"
  echo "MyVar1: $_GARDENVAR_MyVar1"
  echo "MyVar2: $_GARDENVAR_MyVar2"
  echo "IsItAFruit: $_GARDENVAR_IsItAFruit"
  echo "IsApples: $_GARDENVAR_IsApples"
) > "$THIS_SCRIPT_DIR/reap-outputs/${_GARDEN_PLANT_ID}.txt"
