#!/bin/bash

# Set paths to Go and Rust files
MAIN_GO_FILE="main.go"
MAIN_RUST_FILE="init/src/main.rs"
ACTIVE_USER_FILE="activeUser.txt"
IS_RUN_FILE="isRun.txt"

# Path to the lock file that blocks the execution
LOCK_FILE="/tmp/mycmd.lock"

# Check for the existence of the lock file
if [ -e "$LOCK_FILE" ]; then
    echo "Another instance is already running. Exiting."
    exit 1
fi

# Create the lock file
echo "Creating lock file: $LOCK_FILE"
touch "$LOCK_FILE"

# Check if isRun.txt exists and write the appropriate value
if [ -e "$IS_RUN_FILE" ]; then
    echo "true" > "$IS_RUN_FILE"
    echo "File \"$IS_RUN_FILE\" already exists. Written value true."
else
    echo "false" > "$IS_RUN_FILE"
    echo "File \"$IS_RUN_FILE\" created and written value false."
fi

# Check if activeUser.txt exists and delete if it does
if [ -e "$ACTIVE_USER_FILE" ]; then
    echo "Deleting file \"$ACTIVE_USER_FILE\"..."
    rm "$ACTIVE_USER_FILE"
    echo "File deleted."
else
    echo "File \"$ACTIVE_USER_FILE\" does not exist."
fi

# Run the Rust program
rustc "$MAIN_RUST_FILE"
./main

# Run the Go program
go run "$MAIN_GO_FILE"

# Remove the lock file
if [ -e "$LOCK_FILE" ]; then
    rm "$LOCK_FILE"
    echo "Lock file removed."
fi

# Delete isRun.txt after execution
if [ -e "$IS_RUN_FILE" ]; then
    rm "$IS_RUN_FILE"
    echo "File \"$IS_RUN_FILE\" removed."
fi
