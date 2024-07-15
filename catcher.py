import subprocess
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
import traceback

# Адрес электронной почты для отправки сообщений
email_address = "karimovbezan0@gmail.com"
# Логин и пароль для SMTP-сервера (для примера используется SMTP-сервер Gmail)
smtp_server = "smtp.gmail.com"
smtp_port = 587
smtp_user = "karimovbezan0@gmail.com"
smtp_password = "bezhan2009"

def send_email(subject, body):
    try:
        # Создаем сообщение
        msg = MIMEMultipart()
        msg['From'] = smtp_user
        msg['To'] = email_address
        msg['Subject'] = subject

        # Добавляем тело письма
        msg.attach(MIMEText(body, 'plain'))

        # Подключаемся к серверу и отправляем письмо
        server = smtplib.SMTP(smtp_server, smtp_port)
        server.starttls()
        server.login(smtp_user, smtp_password)
        text = msg.as_string()
        server.sendmail(smtp_user, email_address, text)
        server.quit()
    except Exception as e:
        print(f"Failed to send email: {str(e)}")

def run_command(command):
    try:
        # Запускаем команду и ждем завершения
        result = subprocess.run(command, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        print(result.stdout)
    except subprocess.CalledProcessError as e:
        # Ловим ошибки и отправляем их по электронной почте
        error_message = f"Error executing command {command}:\n\n{e.stderr}"
        print(error_message)
        send_email("Error in Go Program", error_message)
    except Exception as e:
        # Ловим любые другие исключения
        error_message = f"An unexpected error occurred while executing command {command}:\n\n{''.join(traceback.format_exception(None, e, e.__traceback__))}"
        print(error_message)
        send_email("Unexpected Error in Go Program", error_message)

if __name__ == "__main__":
    # Команда для выполнения Go-программы
    command = ["go", "run", "main.go"]  # Замените на вашу команду

    run_command(command)
