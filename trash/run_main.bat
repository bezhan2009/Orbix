@echo off
setlocal

rem Set paths to Python files
set MAIN_PYTHON_FILE=init.py
set ACTIVE_USER_FILE=activeUser.txt
set IS_RUN_FILE=isRun.txt
set RUNNING_FILE=running.txt
set CATCHER_FILE=catcher.py

rem Установим переменные окружения только для текущей сессии
set DB_HOST=tiny.db.elephantsql.com
set DB_PORT=5432
set DB_NAME=hzydgvrw
set DB_USER=hzydgvrw
set DB_PASSWORD=7TtuJgOMKm7XVVGL_NheHr4BrpBIMrzz

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

rem Создаем файл running.txt и записываем туда пустую строку
echo. > "%RUNNING_FILE%"

rem Запуск программы на Python
python "%MAIN_PYTHON_FILE%"
python "%CATCHER_FILE%"

rem Дайте небольшую задержку перед удалением файлов
timeout /t 3 /nobreak > NUL

rem Удаляем файлы после завершения программ
del "%RUNNING_FILE%"
del "%IS_RUN_FILE%"

endlocal
