import os
import shutil

# Функция для очистки экрана
def clear_screen():
    if os.name == 'nt':
        os.system('cls')
    else:
        os.system('clear')

# Функция для отображения редактора
def display_editor(filename, lines, cursor):
    clear_screen()
    print(f"==== ORBIX Text Editor ====\n")
    print(f"File: {filename}")
    print("=" * 30)

    # Печать содержимого файла
    for i, line in enumerate(lines):
        marker = ">" if i == cursor else " "
        print(f"{marker} {i + 1}: {line}")

    print("=" * 30)
    print("F2: Save | F3: Quit | F4: Edit Line")

# Чтение содержимого файла
def read_file(filename):
    if os.path.exists(filename):
        with open(filename, 'r') as f:
            return f.readlines()
    return []

# Сохранение файла
def save_file(filename, lines):
    with open(filename, 'w') as f:
        f.writelines(lines)

# Редактирование строки
def edit_line(lines, cursor):
    new_line = input(f"Edit line {cursor + 1}: ")
    lines[cursor] = new_line + '\n'

# Основная функция
def text_editor(filename):
    # Чтение файла
    lines = read_file(filename)
    if not lines:
        lines = ['\n']

    cursor = 0  # Позиция курсора

    while True:
        display_editor(filename, lines, cursor)

        # Получаем ввод пользователя
        key = input("Enter command: ")

        if key == "w":  # Вверх
            cursor = max(0, cursor - 1)
        elif key == "s":  # Вниз
            cursor = min(len(lines) - 1, cursor + 1)
        elif key == "e":  # Редактировать строку
            edit_line(lines, cursor)
        elif key == "f2":  # Сохранить
            save_file(filename, lines)
            print("EditFile saved!")
        elif key == "f3":  # Выход
            print("Exiting editor.")
            break

if __name__ == "__main__":
    filename = input("Enter the filename: ")
    text_editor(filename)
