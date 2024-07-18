@echo off

rem Переходим в директорию exit
cd /d "%~dp0\exit"

rem Строим проект Rust в директории exit
cargo build

rem Переходим в директорию src
cd src

rem Компилируем Rust файл
rustc main.rs -o main.exe

rem Запускаем исполняемый файл main.exe
.\main.exe
