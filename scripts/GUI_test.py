import sys
import subprocess
from PyQt5.QtWidgets import QApplication, QMainWindow, QTabWidget, QWidget, QVBoxLayout, QTextEdit, QPushButton, QHBoxLayout
from PyQt5.QtGui import QFont

class CustomShellApp(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Custom Shell")
        self.setGeometry(100, 100, 800, 600)

        # Настройки шрифта
        self.custom_font = QFont("Consolas", 11)

        # Создаем основной виджет для вкладок
        self.tab_widget = QTabWidget()
        self.setCentralWidget(self.tab_widget)

        # Кнопка для добавления новой вкладки
        add_tab_button = QPushButton("+")
        add_tab_button.clicked.connect(self.add_new_tab)

        # Создаем главный layout и добавляем кнопку
        main_layout = QVBoxLayout()
        main_layout.addWidget(add_tab_button)

        # Создаем виджет для кнопки и layout для главного окна
        central_widget = QWidget()
        central_widget.setLayout(main_layout)

        # Добавляем центральный виджет
        self.setCentralWidget(central_widget)

        # Создаем первую вкладку
        self.create_shell_tab("Tab 1")

    def create_shell_tab(self, tab_name):
        """Создает вкладку с оболочкой."""
        tab_widget = QWidget()
        tab_layout = QVBoxLayout(tab_widget)

        # Создаем текстовое поле для вывода результатов
        output_text = QTextEdit()
        output_text.setReadOnly(True)
        output_text.setFont(self.custom_font)

        # Создаем кнопку для запуска команды
        run_button = QPushButton("Run Command")
        run_button.clicked.connect(lambda: self.run_command(output_text))

        # Добавляем элементы на вкладку
        tab_layout.addWidget(output_text)
        tab_layout.addWidget(run_button)

        # Добавляем вкладку в основной виджет
        self.tab_widget.addTab(tab_widget, tab_name)

        # Запускаем команду go run main.go при создании вкладки
        self.run_command(output_text)

    def run_command(self, output_text):
        """Запускает команду go run main.go и выводит результаты в текстовое поле."""
        try:
            # Запускаем процесс с командой go run main.go
            result = subprocess.run(['go', 'run', 'main.go'], capture_output=True, text=True)

            # Выводим результат выполнения команды
            output_text.append(f"> go run main.go\n")
            output_text.append(result.stdout)
            if result.stderr:
                output_text.append(f"Error: {result.stderr}\n")
        except Exception as e:
            output_text.append(f"Error: {e}\n")

    def add_new_tab(self):
        """Добавляет новую вкладку с оболочкой."""
        tab_count = self.tab_widget.count()
        new_tab_name = f"Tab {tab_count + 1}"
        self.create_shell_tab(new_tab_name)

def main():
    app = QApplication(sys.argv)
    window = CustomShellApp()
    window.show()
    sys.exit(app.exec_())

if __name__ == "__main__":
    main()