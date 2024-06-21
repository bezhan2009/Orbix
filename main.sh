#!/bin/bash

# Определяем путь к файлу main.go
MAIN_GO_FILE="main.go"

# Переходим в корневую директорию проекта
cd "$(dirname "$0")" || exit

# Строим проект Rust в директории init
cd "init" || exit
cargo build

# Переходим в директорию src
cd "src" || exit

# Компилируем Rust файл
rustc main.rs -o main

# Возвращаемся в корневую директорию
cd ../..

# Запускаем Go программу
go run "$MAIN_GO_FILE"
