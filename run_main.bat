@echo off
setlocal

rem Set paths to Go and Rust files
set MAIN_GO_FILE=main.go
set MAIN_RUST_FILE=init\src\main.rs
set ACTIVE_USER_FILE=activeUser.txt
set IS_RUN_FILE=isRun.txt

rem Путь к файлу, который блокирует запуск
set "LOCK_FILE=%TEMP%\mycmd.lock"

rem Проверяем наличие файла блокировки
if exist "%LOCK_FILE%" (
    echo Another instance is already running. Exiting.
    exit /b 1
)

rem Создаем файл блокировки
echo Creating lock file: %LOCK_FILE%
copy NUL "%LOCK_FILE%" > nul

rem Check if isRun.txt exists and write the appropriate value
if exist "%IS_RUN_FILE%" (
    echo true > "%IS_RUN_FILE%"
    echo Файл "%IS_RUN_FILE%" уже существует. Записано значение true.
) else (
    echo false > "%IS_RUN_FILE%"
    echo Файл "%IS_RUN_FILE%" создан и записано значение false.
)

rem Check if activeUser.txt exists and delete if it does
rem Проверяем существование файла activeUser.txt
if exist "%ACTIVE_USER_FILE%" (
    echo Удаляем файл "%ACTIVE_USER_FILE%"...
    del "%ACTIVE_USER_FILE%"
    echo Файл удален.
) else (
    echo Файл "%ACTIVE_USER_FILE%" не существует.
)

rem Запуск программы на Rust
rustc "%MAIN_RUST_FILE%"
.\main.exe

rem Запуск программы на Go
go run "%MAIN_GO_FILE%"

rem Удаляем файл блокировки
if exist "%LOCK_FILE%" (
    del "%LOCK_FILE%"
    echo Файл блокировки удален.
)

rem Delete isRun.txt after execution
if exist "%IS_RUN_FILE%" (
    del "%IS_RUN_FILE%"
    echo Файл "%IS_RUN_FILE%" удален.
)

endlocal
