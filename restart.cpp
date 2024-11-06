#include <iostream>
#include <cstdlib>   // Для функции system
#include <fstream>   // Для проверки наличия файла
#include <chrono>    // Для тайм-аутов
#include <thread>    // Для реализации задержек

bool fileExists(const std::string& filename) {
    std::ifstream file(filename);
    return file.good();
}

int runCommand(const std::string& command) {
    int result = system(command.c_str());
    return result;
}

int main() {
    std::cout << "Restarting Orbix..." << std::endl;
    std::cout << "WARNING: ctrl+C is no longer tracked, so the program will immediately terminate after it" << std::endl;

    bool orbixExists = true;
    bool mainExists = true;

    // Проверка на наличие файла orbix.go в папке scripts
    if (!fileExists("scripts/orbix.go")) {
        orbixExists = false;
        if (!fileExists("orbix.go")) {
            std::cout << "Error: 'orbix.go' not found in the 'scripts' directory or current directory." << std::endl;
            if (!fileExists("main.go")) {
                std::cout << "Error: 'main.go' not found in the current directory." << std::endl;
                mainExists = false;
            }
        }
    }

    if (!orbixExists && !mainExists) {
        return 1;
    }

    const int maxRetries = 3;
    const int retryDelay = 3000; // Задержка между попытками в миллисекундах
    int retries = 0;
    int result = -1;

    while (retries < maxRetries) {
        std::cout << "Attempting to run 'go run orbix.go' in 'scripts' directory (attempt " << retries + 1 << ")..." << std::endl;

        if (orbixExists) {
            // Смена директории на 'scripts' и выполнение команды
            result = runCommand("cd .. && cd scripts && go run orbix.go");
        } else if (mainExists) {
            result = runCommand("cd .. && go run main.go");
        }

        if (result == 0) {
            return 0;  // Команда выполнена успешно
        } else {
            std::cerr << "Error executing command. Return code: " << result << std::endl;
            retries++;
            if (retries < maxRetries) {
                std::cout << "Retrying in " << retryDelay / 1000 << " seconds..." << std::endl;
                std::this_thread::sleep_for(std::chrono::milliseconds(retryDelay));
            } else {
                std::cerr << "Maximum retry attempts reached. Exiting..." << std::endl;
                return result;
            }
        }
    }

    return 0;
}
