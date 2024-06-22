#!/bin/bash

# Укажите полные пути к файлам и директориям
ROOT_DIR="C:/Users/Admin/MyCMD/goCMD"
MAIN_GO_FILE="$ROOT_DIR/main.go"
MAIN_RUST_FILE="$ROOT_DIR/init/src/main.rs"
ACTIVE_USER_FILE="$ROOT_DIR/activeUser.txt"
IS_RUN_FILE="$ROOT_DIR/isRun.txt"

# Переходим в корневую директорию проекта
cd "$ROOT_DIR" || exit

# Создаем и записываем в isRun.txt
echo "true" > "$IS_RUN_FILE"
echo "Файл \"$IS_RUN_FILE\" создан и записано значение true."

# Проверяем существование файла activeUser.txt и удаляем его, если он есть
if [ -f "$ACTIVE_USER_FILE" ]; then
    echo "Удаляем файл \"$ACTIVE_USER_FILE\"..."
    rm "$ACTIVE_USER_FILE"
    echo "Файл удален."
else
    echo "Файл \"$ACTIVE_USER_FILE\" не существует."
fi

# Строим Rust проект
rustc "$MAIN_RUST_FILE"

# Запускаем скомпилированный Rust исполняемый файл
./main.exe

# Пауза (опционально, для просмотра вывода)
read -p "Нажмите Enter для продолжения..."

# Строим Go проект
go run "$MAIN_GO_FILE"
