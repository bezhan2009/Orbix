#!/bin/bash

# Set paths to Python and Rust files
MAIN_PYTHON_FILE="catcher.py"
MAIN_RUST_FILE="init/src/main.rs"
ACTIVE_USER_FILE="activeUser.txt"
IS_RUN_FILE="isRun.txt"

# Check if activeUser.txt exists and delete it if it does
if [ -f "$ACTIVE_USER_FILE" ]; then
    echo "Удаляем файл $ACTIVE_USER_FILE..."
    rm "$ACTIVE_USER_FILE"
    echo "Файл удален."
else
    echo "Файл $ACTIVE_USER_FILE не существует."
fi

# Create isRun.txt and write true to it
echo "true" > "$IS_RUN_FILE"

# Run the Rust program
rustc "$MAIN_RUST_FILE" && ./main

# Run the Go program
python "$MAIN_PYTHON_FILE"

# Delete running.txt
rm running.txt
# Delete isRun.txt
rm isRun.txt
