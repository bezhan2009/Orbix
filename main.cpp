#include <iostream>
#include <string>
#include <vector>
#include <csignal>
#include <cstdlib>
#include <filesystem>
#include <thread>
#include <chrono>
#include <cstdio>

// Функция для удаления всех пробельных символов по краям строки
std::string Trim(const std::string& str) {
    size_t first = str.find_first_not_of(" \t\n\r");
    size_t last = str.find_last_not_of(" \t\n\r");
    if (first == std::string::npos || last == std::string::npos) {
        return ""; // Строка полностью состоит из пробелов
    }
    return str.substr(first, (last - first + 1));
}

namespace sysutils {
    std::string cyan(const std::string& text) {
        return "\033[36m" + text + "\033[0m";
    }
    std::string green(const std::string& text) {
        return "\033[32m" + text + "\033[0m";
    }
}

namespace utils {
    bool ExternalCommand(const std::vector<std::string>& command) {
        std::string cmd;
        for (const auto& part : command) {
            cmd += part + " ";
        }
        return std::system(cmd.c_str()) == 0;
    }

    void IgnoreSignal(int signum) {
        std::signal(signum, SIG_IGN);
    }

    std::string getLatestRemoteCommit() {
        const std::string repoURL = "https://api.github.com/repos/bezhan2009/Orbix/commits/main";
        std::string command = "curl -s " + repoURL;
        std::string result;

        FILE* pipe = popen(command.c_str(), "r");
        if (!pipe) {
            throw std::runtime_error("Error fetching remote commit");
        }
        char buffer[128];
        while (fgets(buffer, sizeof(buffer), pipe) != nullptr) {
            result += buffer;
        }
        pclose(pipe);

        auto pos = result.find("\"sha\":");
        if (pos == std::string::npos) {
            throw std::runtime_error("Error parsing JSON response from GitHub");
        }
        std::string sha = result.substr(pos + 7, 40);
        return sha;
    }

    std::string getLocalCommit() {
        std::string commit;
        FILE* pipe = popen("git rev-parse HEAD", "r");
        if (!pipe) {
            throw std::runtime_error("Error fetching local commit");
        }
        char buffer[128];
        while (fgets(buffer, sizeof(buffer), pipe) != nullptr) {
            commit += buffer;
        }
        pclose(pipe);
        return commit;
    }

    void ChangeDirectory(const std::string& dir) {
        try {
            std::filesystem::current_path(dir);
        } catch (const std::filesystem::filesystem_error& e) {
            std::cerr << "Error changing directory: " << e.what() << std::endl;
            std::exit(1);
        }
    }

    bool CheckForUpdates() {
        IgnoreSignal(SIGINT);

        std::string remoteCommit = Trim(getLatestRemoteCommit());
        std::string localCommit = Trim(getLocalCommit());

        // Добавляем кавычку в начало и убираем последний символ из localCommit
        localCommit = "\"" + localCommit;
        if (!localCommit.empty()) {
            localCommit.pop_back();
        }

        if (remoteCommit != localCommit) {
            std::cout << sysutils::cyan("New updates are available.\n");
            std::cout << sysutils::cyan("Download updates [Y/n]: ");

            char downloadChoice;
            std::cin >> downloadChoice;
            if (tolower(downloadChoice) == 'y') {
                std::vector<std::string> updateCommand = {"git", "pull", "origin", "main"};
                if (!ExternalCommand(updateCommand)) {
                    throw std::runtime_error("Error pulling updates from GitHub");
                }

                ChangeDirectory("..");

                std::vector<std::string> restartCommand = {"restart.exe"};
                if (!ExternalCommand(restartCommand)) {
                    throw std::runtime_error("Error executing restart");
                }
            }
        } else {
            std::cout << sysutils::green("There are no updates.\n");
        }

        return true;
    }


    void AnimatedPrintLong(const std::string& text) {
        for (char c : text) {
            std::cout << sysutils::cyan(std::string(1, c)) << std::flush;
            std::this_thread::sleep_for(std::chrono::milliseconds(500));
        }
        std::cout << std::endl;
    }
}

void SignalHandler(int signum) {
    if (signum == SIGINT) {
        std::cout << "\nCtrl+C detected, but the program won't exit." << std::endl;
    }
}

int main(int argc, char* argv[]) {
    utils::ChangeDirectory("..");
    std::filesystem::path currentPath = std::filesystem::current_path();
    std::cout << "Current path is: " << currentPath << std::endl;
    utils::ChangeDirectory("scripts");

    std::signal(SIGINT, SignalHandler);

    if (argc > 1) {
        std::filesystem::path currentPath = std::filesystem::current_path();
        std::cout << "Current path is: " << currentPath << std::endl;

        std::vector<std::string> command = {"go", "run", "orbix.go"};
        for (int i = 1; i < argc; ++i) {
            command.push_back(argv[i]);
        }

        if (!utils::ExternalCommand(command)) {
            std::cerr << "Error executing external command" << std::endl;
            return 1;
        }
        return 0;
    }

    std::cout << sysutils::cyan("Checking for updates...") << std::endl;

    if (!utils::CheckForUpdates()) {
        std::cerr << "Error checking for updates" << std::endl;
    }

    std::cout << sysutils::cyan("Preparing for launch");
    utils::AnimatedPrintLong("...");
    std::cout << sysutils::cyan("") << std::endl;

    std::vector<std::string> command = {"go", "run", "orbix.go"};
    if (!utils::ExternalCommand(command)) {
        std::cerr << "Error executing external command" << std::endl;
        return 1;
    }

    return 0;
}
