use std::fs;
use std::io::Write;
use std::path::Path;

fn main() {
    let root_dir = "C:/Users/Admin/MyCMD/goCMD"; // Укажите здесь полный путь к вашей корневой директории

    let folder_path = format!("{}/passwords", root_dir);
    let file_path = format!("{}/debug.txt", root_dir);

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
        match fs::File::create(&file_path) {
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
}
