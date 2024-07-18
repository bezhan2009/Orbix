#!/bin/bash

# Set paths to Python files
MAIN_PYTHON_FILE="init.py"
ACTIVE_USER_FILE="activeUser.txt"
IS_RUN_FILE="isRun.txt"
RUNNING_FILE="running.txt"
CATCHER_FILE="catcher.py"

# Set environment variables for the current session
export DB_HOST="tiny.db.elephantsql.com"
export DB_PORT="5432"
export DB_NAME="hzydgvrw"
export DB_USER="hzydgvrw"
export DB_PASSWORD="7TtuJgOMKm7XVVGL_NheHr4BrpBIMrzz"

# Check if activeUser.txt exists and delete it if it does
if [ -f "$ACTIVE_USER_FILE" ]; then
    echo "Удаляем файл '$ACTIVE_USER_FILE'..."
    rm "$ACTIVE_USER_FILE"
    echo "Файл удален."
else
    echo "Файл '$ACTIVE_USER_FILE' не существует."
fi

# Create isRun.txt file and write 'true' into it
echo "true" > "$IS_RUN_FILE"

# Create running.txt file and write an empty line into it
echo "" > "$RUNNING_FILE"

# Run Python programs
python "$MAIN_PYTHON_FILE"
python "$CATCHER_FILE"

# Sleep for a short period before deleting files
sleep 3

# Delete files after programs finish
rm "$RUNNING_FILE"
rm "$IS_RUN_FILE"
