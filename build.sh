#!/bin/bash
cmake -B build
cmake --build build
# shellcheck disable=SC2164
cd build
./main.exe
