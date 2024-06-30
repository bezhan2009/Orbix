#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <libpq-fe.h>
#include <curl/curl.h>
#include <unistd.h>

void show_help() {
    printf("Orbix Package Manager (OPM)\n");
    printf("Usage: opm [command] [options]\n");
    printf("Commands:\n");
    printf("  install <package>   Install a package\n");
    printf("  remove <package>    Remove a package\n");
    printf("  update              Update the package manager\n");
    printf("  list                List installed packages\n");
    printf("  help                Show this help message\n");
}

void do_exit(PGconn *conn) {
    PQfinish(conn);
    exit(1);
}

size_t write_data(void *ptr, size_t size, size_t nmemb, FILE *stream) {
    size_t written = fwrite(ptr, size, nmemb, stream);
    return written;
}

void download_file(const char *url, const char *filename) {
    CURL *curl;
    FILE *fp;
    CURLcode res;

    curl = curl_easy_init();
    if (curl) {
        fp = fopen(filename, "wb");
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_data);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, fp);
        res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            fprintf(stderr, "Download failed: %s\n", curl_easy_strerror(res));
        }
        fclose(fp);
        curl_easy_cleanup(curl);
    }
}

void install_package(const char *package_name) {
    const char *conninfo = "dbname=yourdbname user=youruser password=yourpassword hostaddr=yourhostaddr port=yourport";
    PGconn *conn = PQconnectdb(conninfo);

    if (PQstatus(conn) == CONNECTION_BAD) {
        fprintf(stderr, "Connection to database failed: %s\n", PQerrorMessage(conn));
        do_exit(conn);
    }

    char query[256];
    snprintf(query, sizeof(query), "SELECT url FROM Packages WHERE name='%s'", package_name);

    PGresult *res = PQexec(conn, query);

    if (PQresultStatus(res) != PGRES_TUPLES_OK) {
        fprintf(stderr, "No data retrieved: %s\n", PQerrorMessage(conn));
        PQclear(res);
        do_exit(conn);
    }

    if (PQntuples(res) > 0) {
        const char *url = PQgetvalue(res, 0, 0);
        printf("Downloading package from: %s\n", url);
        download_file(url, package_name);
        printf("Package '%s' installed successfully.\n", package_name);
    } else {
        printf("Package '%s' not found.\n", package_name);
    }

    PQclear(res);
    PQfinish(conn);
}

void remove_package(const char *package_name) {
    char filename[FILENAME_MAX] = "./";
    strcat(filename, package_name);
    strcat(filename, ".pkg");

    if (remove(filename) == 0) {
        printf("Package '%s' removed successfully.\n", package_name);
    } else {
        printf("Error: Could not remove package '%s'.\n", package_name);
    }
}

void update_package_manager() {
    printf("Updating package manager...\n");
    const char *update_url = "your_update_url"; // URL для обновления пакетного менеджера
    const char *update_filename = "opm_new";

    download_file(update_url, update_filename);

    // Замена текущего исполняемого файла новым
    if (rename(update_filename, "opm") != 0) {
        fprintf(stderr, "Failed to update package manager.\n");
    } else {
        printf("Package manager updated successfully.\n");
    }
}

void list_installed_packages() {
    printf("Listing installed packages:\n");

    // Для простоты предположим, что установленные пакеты находятся в текущем каталоге с расширением .pkg
    system("ls *.pkg");
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        printf("Error: No command provided.\n");
        show_help();
        return 1;
    }

    if (strcmp(argv[1], "install") == 0) {
        if (argc < 3) {
            printf("Error: No package name provided.\n");
            return 1;
        }
        const char *package_name = argv[2];
        install_package(package_name);
    } else if (strcmp(argv[1], "remove") == 0) {
        if (argc < 3) {
            printf("Error: No package name provided.\n");
            return 1;
        }
        const char *package_name = argv[2];
        remove_package(package_name);
    } else if (strcmp(argv[1], "update") == 0) {
        update_package_manager();
    } else if (strcmp(argv[1], "list") == 0) {
        list_installed_packages();
    } else if (strcmp(argv[1], "help") == 0) {
        show_help();
    } else {
        printf("Error: Unknown command '%s'.\n", argv[1]);
        show_help();
    }

    return 0;
}
