use std::fs;
use std::fs::File;
use std::io::{Read, Write};
use std::path::Path;

fn main() {
    println!("initializer started!");

    let root_dir = "C:/Users/Admin/MyCMD/goCMD"; // Укажите здесь полный путь к вашей корневой директории

    let folder_path = format!("{}/passwords", root_dir);
    let file_path = format!("{}/debug.txt", root_dir);
    let file_run_path = format!("{}/running.txt", root_dir);
    let file_user_path = format!("{}/activeUser.txt", root_dir);

    // Проверяем существование папки, если нет - создаем
    if !Path::new(&folder_path).exists() {
        match fs::create_dir(&folder_path) {
            Ok(_) => println!("Папка '{}' создана.", folder_path),
            Err(e) => eprintln!("Ошибка при создании папки '{}': {}", folder_path, e),
        }
    } else {
        println!("Папка '{}' уже существует.", folder_path);
    }

    // Проверяем существование файла, если нет - создаем
    if !Path::new(&file_path).exists() {
        match File::create(&file_path) {
            Ok(mut file) => {
                println!("Файл '{}' создан.", file_path);
                if let Err(e) = file.write_all(b"Debug information") {
                    eprintln!("Ошибка при записи в файл '{}': {}", file_path, e);
                }
            }
            Err(e) => eprintln!("Ошибка при создании файла '{}': {}", file_path, e),
        }
    } else {
        println!("Файл '{}' уже существует.", file_path);
    }

    // Проверяем существование файла activeUser.txt и выходим с паникой, если он существует
    if Path::new(&file_user_path).exists() {
        panic!("Файл '{}' уже существует. Программа завершена.", file_user_path);
    }

    // Проверяем существование файла running.txt и выходим с паникой, если он существует
    if Path::new(&file_run_path).exists() {
        panic!("Файл '{}' уже существует. Программа завершена.", file_run_path);
    }

    // Считываем данные из файла activeUser.txt
    let mut user_file = File::open(&file_user_path)
        .expect(&format!("Не удалось открыть файл '{}'", file_user_path));
    let mut user_data = String::new();
    user_file.read_to_string(&mut user_data)
        .expect(&format!("Не удалось прочитать данные из файла '{}'", file_user_path));

    // Создаем файл running.txt и записываем данные из файла activeUser.txt
    match File::create(&file_run_path) {
        Ok(mut file) => {
            println!("Файл '{}' создан.", file_run_path);
            if let Err(e) = file.write_all(user_data.as_bytes()) {
                eprintln!("Ошибка при записи в файл '{}': {}", file_run_path, e);
            }
        }
        Err(e) => eprintln!("Ошибка при создании файла '{}': {}", file_run_path, e),
    }

    println!("initializer completed!");
}
