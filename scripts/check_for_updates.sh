#!/bin/bash
# shellcheck disable=SC2164
cd ..

cmake -B build
cmake --build build

# shellcheck disable=SC2164
cd build
./main.exe
