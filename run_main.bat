@echo off

rem Set paths to Go and Rust files
set MAIN_GO_FILE=main.go
set MAIN_RUST_FILE=init\src\main.rs

rem Change directory to the root of the project
cd /d "%~dp0"

rem Build Rust project
rustc "%MAIN_RUST_FILE%"

rem Run the compiled Rust executable
.\main.exe

rem Pause to see output (optional)
pause

rem Build Go project
go run "%MAIN_GO_FILE%"

