#!/bin/bash

# Определяем путь к файлу main.go
MAIN_GO_FILE="main.go"

# Переходим в директорию init
cd "init" || exit

# Строим проект Go
cargo build

# Переходим в директорию src
cd "src" || exit

# Компилируем Rust файл
rustc main.rs

# Переходим обратно в исходную директорию
cd ..

# Запускаем Go программу
go run "$MAIN_GO_FILE"
