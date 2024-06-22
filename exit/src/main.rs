use std::process::{Command, exit};
use std::fs;

fn main() {
    let file_path = "C:/Users/Admin/MyCMD/goCMD/activeUser.txt";

    println!("Executing exit command...");

    // Проверяем существование файла activeUser.txt
    if fs::metadata(&file_path).is_ok() {
        // Удаляем файл activeUser.txt
        if let Err(e) = fs::remove_file(&file_path) {
            eprintln!("Failed to delete file '{}': {}", file_path, e);
        } else {
            println!("File '{}' deleted successfully.", file_path);
        }
    } else {
        println!("File '{}' does not exist.", file_path);
    }

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
