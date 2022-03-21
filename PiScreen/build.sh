#!/bin/bash

mkdir -p build
cd build

cmake .. -GNinja -DCMAKE_BUILD_TYPE=Release
cmake --build . --parallel
