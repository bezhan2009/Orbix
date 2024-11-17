@echo off
cd ..
cmake -B build
cmake --build build
cd build
main.exe
pause
