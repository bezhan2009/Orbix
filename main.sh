#!/bin/bash

# Пути к файлам Go и Rust
MAIN_GO_FILE="main.go"
MAIN_RUST_FILE="init/src/main.rs"
ACTIVE_USER_FILE="activeUser.txt"
RUNNING_FILE="running.txt"

# Удаляем файлы activeUser.txt и running.txt, если они существуют
if [ -f "$ACTIVE_USER_FILE" ]; then
    echo "Удаляем файл \"$ACTIVE_USER_FILE\"..."
    rm "$ACTIVE_USER_FILE"
    echo "Файл \"$ACTIVE_USER_FILE\" удален."
else
    echo "Файл \"$ACTIVE_USER_FILE\" не существует."
fi

if [ -f "$RUNNING_FILE" ]; then
    echo "Удаляем файл \"$RUNNING_FILE\" после выполнения программы на Go."
fi

# Компилируем и запускаем программу на Rust
rustc "$MAIN_RUST_FILE"
./main.exe

# Запускаем программу на Go и удаляем файл running.txt после ее завершения
go run "$MAIN_GO_FILE"
rm -f "$RUNNING_FILE"

echo "Завершено удаление файла \"$RUNNING_FILE\"."
