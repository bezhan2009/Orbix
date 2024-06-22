import platform
import subprocess
import os

def run_exit_script():
    current_dir = os.path.dirname(os.path.realpath(__file__))

    if platform.system() == 'Windows':
        script_path = os.path.join(current_dir, 'exit.bat')
        subprocess.call(script_path, shell=True)
    elif platform.system() == 'Linux' or platform.system() == 'Darwin':  # Unix or MacOS
        script_path = os.path.join(current_dir, 'exit.sh')
        subprocess.call(['bash', script_path])
    else:
        print("Unsupported operating system")

run_exit_script()
