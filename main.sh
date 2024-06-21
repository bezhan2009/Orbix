#!/bin/bash

MAIN_GO_FILE="main.go"

# shellcheck disable=SC2164
cd "$(dirname "$MAIN_GO_FILE")"


go run "$(basename "$MAIN_GO_FILE")"
