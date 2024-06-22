@echo off

rem Set paths to Go and Rust files
set MAIN_GO_FILE=main.go
set MAIN_RUST_FILE=init\src\main.rs
set ACTIVE_USER_FILE=activeUser.txt

rem Change directory to the root of the project
cd /d "%~dp0"

rem Check if activeUser.txt exists and exit if it does
if exist "%ACTIVE_USER_FILE%" (
    echo Файл "%ACTIVE_USER_FILE%" существует. Программа завершена.
    exit /b 1
)

rem Build Rust project
rustc "%MAIN_RUST_FILE%"

rem Run the compiled Rust executable
.\main.exe

rem Pause to see output (optional)
pause

rem Build Go project
go run "%MAIN_GO_FILE%"
