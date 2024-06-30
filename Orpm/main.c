#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <libpq-fe.h>

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

void install_package(const char *package_name) {
    const char *conninfo = "dbname=yourdbname user=youruser password=yourpassword hostaddr=yourhostaddr port=yourport";
    PGconn *conn = PQconnectdb(conninfo);

    if (PQstatus(conn) == CONNECTION_BAD) {
        fprintf(stderr, "Connection to database failed: %s\n", PQerrorMessage(conn));
        do_exit(conn);
    }

    char query[256];
    snprintf(query, sizeof(query), "SELECT * FROM packages WHERE name='%s'", package_name);

    PGresult *res = PQexec(conn, query);

    if (PQresultStatus(res) != PGRES_TUPLES_OK) {
        fprintf(stderr, "No data retrieved: %s\n", PQerrorMessage(conn));
        PQclear(res);
        do_exit(conn);
    }

    // Обработка данных пакета
    // Например:
    int nFields = PQnfields(res);
    for (int i = 0; i < PQntuples(res); i++) {
        for (int j = 0; j < nFields; j++) {
            printf("%s: %s\n", PQfname(res, j), PQgetvalue(res, i, j));
        }
    }

    PQclear(res);
    PQfinish(conn);
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
        printf("Removing package...\n");
    } else if (strcmp(argv[1], "update") == 0) {
        printf("Updating package manager...\n");
    } else if (strcmp(argv[1], "list") == 0) {
        printf("Listing installed packages...\n");
    } else if (strcmp(argv[1], "help") == 0) {
        show_help();
    } else {
        printf("Error: Unknown command '%s'.\n", argv[1]);
        show_help();
    }

    return 0;
}
