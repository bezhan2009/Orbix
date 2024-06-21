use std::process::{Command, exit};

fn main() {
    println!("Executing exit command...");

    // Выполнение команды `cmd /C exit` для завершения командной оболочки
    let status = Command::new("cmd")
        .arg("/C")
        .arg("exit")
        .status()
        .expect("Failed to execute exit command");

    // Проверяем статус завершения команды
    if !status.success() {
        println!("Failed to execute exit command");
    }

    // Завершаем выполнение программы Rust с кодом 1 (необязательно)
    exit(1);
}
