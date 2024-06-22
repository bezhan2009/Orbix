@echo off

rem Set paths to Go and Rust files
set MAIN_GO_FILE=main.go
set MAIN_RUST_FILE=init\src\main.rs
set ACTIVE_USER_FILE=activeUser.txt
set IS_RUN_FILE=isRun.txt

rem Change directory to the root of the project
cd /d "%~dp0"

rem Create and write to isRun.txt
echo true > "%IS_RUN_FILE%"
echo Файл "%IS_RUN_FILE%" создан и записано значение true.

rem Check if activeUser.txt exists and exit if it does
rem Проверяем существование файла activeUser.txt
if exist "%ACTIVE_USER_FILE%" (
    echo Удаляем файл "%ACTIVE_USER_FILE%"...
    del "%ACTIVE_USER_FILE%"
    echo Файл удален.
) else (
    echo Файл "%ACTIVE_USER_FILE%" не существует.
)

rem Build Rust project
rustc "%MAIN_RUST_FILE%"

rem Run the compiled Rust executable
.\main.exe

rem Pause to see output (optional)
pause

rem Build Go project
go run "%MAIN_GO_FILE%"
