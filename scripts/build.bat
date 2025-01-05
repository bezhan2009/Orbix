@echo off

:: Create the 'bin' folder if it does not exist
if not exist "bin" (
    mkdir bin
)

:: Compile the Go file and save it in the 'bin' folder
go build -gcflags=all="-B" -ldflags="-s -w" -o bin\orbix.exe orbix.go

:: Add the 'bin' folder to the PATH environment variable
setx PATH "%PATH%;%CD%\bin"

:: Run the executable file
bin\orbix.exe

pause
