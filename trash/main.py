import platform
import subprocess
import os
import sys

def run():
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
            print(f"File not found: {script_path}")
        else:
            subprocess.call(script_path, shell=True)
    elif platform.system() == 'Linux' or platform.system() == 'Darwin':  # Unix or MacOS
        script_path = os.path.join(current_dir, 'main.sh')
        print(f"Script path: {script_path}")

        if not os.path.isfile(script_path):
            print(f"File not found: {script_path}")
        else:
            subprocess.call(['bash', script_path])
    else:
        print("Unsupported operating system")

run()
