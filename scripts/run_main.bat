@echo off
setlocal

rem Set paths to Go and Rust files
set MAIN_GO_FILE=main.go
set MAIN_RUST_FILE=init\src\main.rs
set ACTIVE_USER_FILE=activeUser.txt
set IS_RUN_FILE=isRun.txt

rem Проверяем наличие файла activeUser.txt и удаляем его, если он существует
if exist "%ACTIVE_USER_FILE%" (
    echo Удаляем файл "%ACTIVE_USER_FILE%"...
    del "%ACTIVE_USER_FILE%"
    echo Файл удален.
) else (
    echo Файл "%ACTIVE_USER_FILE%" не существует.
)

rem Создаем файл isRun.txt и записываем туда true
echo true > "%IS_RUN_FILE%"

rem Запуск программы на Rust
rustc "%MAIN_RUST_FILE%"
.\main.exe

rem Запуск программы на Go
go run "%MAIN_GO_FILE%"

del running.txt
del isRun.txt

endlocal
