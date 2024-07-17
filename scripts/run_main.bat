@echo off
setlocal

rem Set paths to Python files
set MAIN_PYTHON_FILE=init.py
set ACTIVE_USER_FILE=activeUser.txt
set IS_RUN_FILE=isRun.txt
set RUNNING_FILE=running.txt
set CATCHER_FILE=catcher.py

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

rem Создаем файл running.txt и записываем туда нечего
echo "" > "%RUNNING_FILE%"

rem Запуск программы на Python
python "%MAIN_PYTHON_FILE%"
python "%CATCHER_FILE%"

del running.txt
del isRun.txt

endlocal
