import platform
import subprocess
import os
import sys
import traceback
import psycopg2
from dotenv import load_dotenv

# Загрузка переменных окружения из файла .env
load_dotenv()

# Переменные для подключения к PostgreSQL
db_host = os.getenv("DB_HOST")
db_port = os.getenv("DB_PORT")
db_name = os.getenv("DB_NAME")
db_user = os.getenv("DB_USER")
db_password = os.getenv("DB_PASSWORD")

def insert_error_to_db(command, error_message):
    try:
        # Подключение к базе данных
        conn = psycopg2.connect(
            dbname=db_name,
            user=db_user,
            password=db_password,
            host=db_host,
            port=db_port
        )
        cursor = conn.cursor()

        # Создание таблицы, если она еще не существует
        cursor.execute("""
            CREATE TABLE IF NOT EXISTS error_log (
                id SERIAL PRIMARY KEY,
                command TEXT,
                error_message TEXT,
                timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            )
        """)

        # Вставка данных об ошибке
        cursor.execute("""
            INSERT INTO error_log (command, error_message)
            VALUES (%s, %s)
        """, (command, error_message))

        # Подтверждение изменений и закрытие соединения
        conn.commit()
        cursor.close()
        conn.close()
    except Exception as e:
        print(f"Failed to insert error into database: {e}")

def run_go_script():
    try:
        if getattr(sys, 'frozen', False):  # Проверка, если скрипт запущен как .exe
            current_dir = os.path.dirname(sys.executable)
        else:
            current_dir = os.path.dirname(os.path.realpath(__file__))

        print(f"Current directory: {current_dir}")
        os.chdir(current_dir)  # Изменение текущей директории на директорию скрипта

        # Проверка существования файла main.go в текущей директории
        go_file = os.path.join(current_dir, "main.go")
        if not os.path.exists(go_file):
            raise FileNotFoundError(f"{go_file} does not exist")

        # Запуск Go файла main.go
        go_command = "go run main.go"
        print(f"Running command: {go_command}")

        with open('stdout.log', 'w') as stdout_log, open('stderr.log', 'w') as stderr_log:
            process = subprocess.Popen(go_command, shell=True, stdout=stdout_log, stderr=stderr_log, text=True)

            try:
                process.wait(timeout=60)  # Увеличен тайм-аут до 60 секунд
            except subprocess.TimeoutExpired:
                process.kill()
                raise subprocess.TimeoutExpired(go_command, 60)

            if process.returncode != 0:
                with open('stderr.log', 'r') as stderr_log:
                    stderr = stderr_log.read()
                raise subprocess.CalledProcessError(process.returncode, go_command, stderr=stderr)
            else:
                with open('stdout.log', 'r') as stdout_log:
                    stdout = stdout_log.read()
                print("Go script executed successfully.")
                print(f"stdout: {stdout}")

    except subprocess.TimeoutExpired as e:
        error_message = f"Command '{e.cmd}' timed out after {e.timeout} seconds"
        print(error_message)
        insert_error_to_db(go_command, error_message)
    except Exception as e:
        error_message = f"An error occurred:\n\n{''.join(traceback.format_exception(None, e, e.__traceback__))}"
        print(error_message)
        insert_error_to_db(go_command, error_message)  # Передаем команду, которую запускали, и сообщение об ошибке

def run():
    try:
        if platform.system() == 'Windows':
            run_go_script()
        else:
            raise OSError("Unsupported operating system")
    except Exception as e:
        error_message = f"An error occurred:\n\n{''.join(traceback.format_exception(None, e, e.__traceback__))}"
        print(error_message)
        insert_error_to_db(" ".join(sys.argv), error_message)  # Передаем команду, которую запускали, и сообщение об ошибке

if __name__ == "__main__":
    run()
