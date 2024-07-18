import os
import subprocess

def main():
    file_path = "C:/Users/Admin/MyCMD/goCMD/activeUser.txt"

    print("Executing exit command...")

    # Проверяем существование файла activeUser.txt
    if os.path.exists(file_path):
        try:
            os.remove(file_path)
            print(f"File '{file_path}' deleted successfully.")
        except Exception as e:
            print(f"Failed to delete file '{file_path}': {e}")
    else:
        print(f"File '{file_path}' does not exist.")

    # Выполнение команды `cmd /C exit` для завершения командной оболочки
    try:
        result = subprocess.run(["cmd", "/C", "exit"], check=True)
        if result.returncode != 0:
            print("Failed to execute exit command")
    except subprocess.CalledProcessError:
        print("Failed to execute exit command")

if __name__ == "__main__":
    main()
