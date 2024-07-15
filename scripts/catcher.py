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

def run():
    try:
        if getattr(sys, 'frozen', False):  # Проверка, если скрипт запущен как .exe
            current_dir = os.path.dirname(sys.executable)
        else:
            current_dir = os.path.dirname(os.path.realpath(__file__))

        print(f"Current directory: {current_dir}")
        os.chdir(current_dir)  # Изменение текущей директории на директорию скрипта

        if platform.system() == 'Windows':
            script_path = os.path.join(current_dir, 'run_main.bat')
            print(f"Script path: {script_path}")

            if not os.path.isfile(script_path):
                raise FileNotFoundError(f"File not found: {script_path}")
            else:
                subprocess.call(script_path, shell=True)
        elif platform.system() == 'Linux' or platform.system() == 'Darwin':  # Unix or MacOS
            script_path = os.path.join(current_dir, 'main.sh')
            print(f"Script path: {script_path}")

            if not os.path.isfile(script_path):
                raise FileNotFoundError(f"File not found: {script_path}")
            else:
                subprocess.call(['bash', script_path])
        else:
            raise OSError("Unsupported operating system")
    except Exception as e:
        error_message = f"An error occurred:\n\n{''.join(traceback.format_exception(None, e, e.__traceback__))}"
        print(error_message)
        insert_error_to_db(" ".join(sys.argv), error_message)  # Передаем команду, которую запускали, и сообщение об ошибке

if __name__ == "__main__":
    run()
