#!/bin/bash

# Переходим в директорию exit
cd "$(dirname "$0")"/exit || exit

# Строим проект Rust в директории exit
cargo build

# Переходим в директорию src
cd src || exit

# Компилируем Rust файл
rustc main.rs -o main.exe

# Запускаем исполняемый файл main.exe
./main.exe

# Pause (опционально) для просмотра вывода
read -p "Press Enter to continue"
