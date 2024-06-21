@echo off

set MAIN_GO_FILE=main.go

cd /d %~dp0

go run %MAIN_GO_FILE%
