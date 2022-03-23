#!/bin/bash

export BaseDir=$(realpath $(dirname .))

for x in ./*/build.sh; do
    export d=$(realpath $(dirname $(pwd)/$x))
    export f=$(realpath $(pwd)/$x)
    cd $d && echo "cd to $d"
    $f
    cd $BaseDir
done
