import os
import tkinter as tk
from tkinter import ttk, font

class CustomShellApp:
    def __init__(self, root):
        self.root = root
        self.root.title("Custom Shell")
        self.root.configure(bg='#000000')  # Черный фон

        # Настройки шрифта
        self.custom_font = font.Font(family="Consolas", size=11)

        # Создаем основной контейнер для вкладок
        self.notebook = ttk.Notebook(self.root)
        self.notebook.pack(expand=True, fill=tk.BOTH, padx=10, pady=10)

        # Кнопка для добавления новой вкладки
        add_tab_button = ttk.Button(self.root, text="+", command=self.add_new_tab)
        add_tab_button.pack(side=tk.TOP, padx=10, pady=10)

        # Создаем первую вкладку
        self.create_shell_tab("Tab 1")

    def create_shell_tab(self, tab_name):
        """Создает вкладку с оболочкой."""
        tab_frame = ttk.Frame(self.notebook)
        tab_frame.pack(fill=tk.BOTH, expand=True)

        # Создаем текстовое поле для вывода результатов
        output_text = tk.Text(tab_frame, wrap=tk.WORD, height=20, width=80, bg='#000000', fg='#ffffff', insertbackground='#ffffff', font=self.custom_font)
        output_text.pack(padx=10, pady=10)

        # Функция для запуска CMD
        def run_cmd():
            cmd_output = os.popen('cmd').read()
            output_text.insert(tk.END, cmd_output)

        # Запускаем CMD при создании вкладки
        run_cmd()

        # Добавляем вкладку в основной ноутбук
        self.notebook.add(tab_frame, text=tab_name)

    def add_new_tab(self):
        """Добавляет новую вкладку с оболочкой."""
        tab_count = self.notebook.index(tk.END)
        new_tab_name = f"Tab {tab_count + 1}"
        self.create_shell_tab(new_tab_name)

def main():
    root = tk.Tk()
    app = CustomShellApp(root)
    root.mainloop()

if __name__ == "__main__":
    main()
