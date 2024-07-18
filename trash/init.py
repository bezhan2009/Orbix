import os

def main():
    print("initializer started!")

    # Получаем текущую директорию скрипта
    script_dir = os.path.dirname(os.path.abspath(__file__))

    # Определяем корневую директорию на одну папку выше
    root_dir = os.path.join(script_dir, '..')

    # Нормализуем путь, чтобы получить абсолютный путь
    root_dir = os.path.abspath(root_dir)

    folder_path = os.path.join(root_dir, "passwords")
    file_path = os.path.join(root_dir, "debug.txt")
    file_run_path = os.path.join(root_dir, "running.txt")
    file_user_path = os.path.join(root_dir, "activeUser.txt")

    print("Script directory:", script_dir)
    print("Root directory:", root_dir)
    print("Folder path:", folder_path)
    print("Debug file path:", file_path)
    print("Running file path:", file_run_path)
    print("Active user file path:", file_user_path)

    # Проверяем существование папки, если нет - создаем
    if not os.path.exists(folder_path):
        try:
            os.makedirs(folder_path)
            print(f"Папка '{folder_path}' создана.")
        except Exception as e:
            print(f"Ошибка при создании папки '{folder_path}': {e}")
    else:
        print(f"Папка '{folder_path}' уже существует.")

    # Проверяем существование файла, если нет - создаем
    if not os.path.exists(file_path):
        try:
            with open(file_path, 'w') as file:
                print(f"Файл '{file_path}' создан.")
                file.write("Debug information")
        except Exception as e:
            print(f"Ошибка при создании файла '{file_path}': {e}")
    else:
        print(f"Файл '{file_path}' уже существует.")

    # Проверяем существование файла activeUser.txt и выходим с паникой, если он существует
    if os.path.exists(file_user_path):
        print(f"Файл '{file_user_path}' уже существует. Программа завершена.")
        return

    # Проверяем существование файла running.txt и выходим с паникой, если он существует
    if os.path.exists(file_run_path):
        print(f"Файл '{file_run_path}' уже существует. Программа завершена.")
        return

    # Считываем данные из файла activeUser.txt
    try:
        with open(file_user_path, 'r') as user_file:
            user_data = user_file.read()
    except Exception as e:
        raise Exception(f"Не удалось открыть или прочитать файл '{file_user_path}': {e}")

    # Создаем файл running.txt и записываем данные из файла activeUser.txt
    try:
        with open(file_run_path, 'w') as file:
            print(f"Файл '{file_run_path}' создан.")
            file.write(user_data)
    except Exception as e:
        print(f"Ошибка при создании файла '{file_run_path}': {e}")

    print("initializer completed!")

if __name__ == "__main__":
    main()
