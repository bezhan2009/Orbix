#!/bin/bash

# Set paths to Python files
MAIN_PYTHON_FILE="init.py"
ACTIVE_USER_FILE="activeUser.txt"
IS_RUN_FILE="isRun.txt"
RUNNING_FILE="running.txt"
CATCHER_FILE="catcher.py"

# Проверяем наличие файла activeUser.txt и удаляем его, если он существует
if [ -f "$ACTIVE_USER_FILE" ]; then
    echo "Удаляем файл '$ACTIVE_USER_FILE'..."
    rm "$ACTIVE_USER_FILE"
    echo "Файл удален."
else
    echo "Файл '$ACTIVE_USER_FILE' не существует."
fi

# Создаем файл isRun.txt и записываем туда true
echo "true" > "$IS_RUN_FILE"

# Создаем файл running.txt и записываем туда нечего
echo "" > "$RUNNING_FILE"

# Запуск программы на Python
python "$MAIN_PYTHON_FILE"
python "$CATCHER_FILE"

# Удаляем временные файлы
rm running.txt
rm isRun.txt
