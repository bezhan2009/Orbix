#include <iostream>
#include <cstdlib>  // Для функции system
#include <fstream>  // Для проверки наличия файла
#include <chrono>   // Для тайм-аутов
#include <thread>   // Для реализации задержек

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

    // Проверка на наличие файла orbix.go в папке scripts
    if (!fileExists("scripts/orbix.go")) {
        std::cerr << "Error: 'orbix.go' not found in the 'scripts' directory." << std::endl;
        return 1;
    }

    const int maxRetries = 3;
    const int retryDelay = 3000; // Задержка между попытками в миллисекундах
    int retries = 0;

    while (retries < maxRetries) {
        std::cout << "Attempting to run 'go run orbix.go' in 'scripts' directory (attempt " << retries + 1 << ")..." << std::endl;

        // Смена директории на 'scripts' и выполнение команды
        int result = runCommand("cd scripts && go run orbix.go");

        if (result == 0) {
            return 0;
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
