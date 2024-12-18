@echo off

go build -gcflags=all="-B" -ldflags="-s -w" orbix.go

orbix.exe

pause
