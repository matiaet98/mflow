#!/bin/bash

if [ $# -ne 2 ]
  then
    echo "Usage: mflow.sh <configfile_path> <taskfile_path>"
    echo "Ej: mflow.sh config.json tasks1.json"
    exit 1
fi

source /home/hadoop/.bash_profile

/usr/local/fiscar/mflow/mflow --config $1 --taskfile $2
