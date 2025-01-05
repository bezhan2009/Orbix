#!/bin/bash

# Create the 'bin' folder if it does not exist
mkdir -p bin

# Compile the Go file and save it in the 'bin' folder
go build -gcflags=all="-B" -ldflags="-s -w" -o bin/orbix orbix.go

# Add the 'bin' folder to the PATH environment variable
# shellcheck disable=SC2155
export PATH="$PATH:$(pwd)/bin"
# shellcheck disable=SC2016
echo 'export PATH="$PATH:$(pwd)/bin"' >> ~/.bashrc

# Run the executable file
bin/orbix
