package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func showHelp() {
	fmt.Println("Orbix Package Manager (ORPM)")
	fmt.Println("Usage: opm [command] [options]")
	fmt.Println("Commands:")
	fmt.Println("  install <package>   Install a package")
	fmt.Println("  remove <package>    Remove a package")
	fmt.Println("  update              Update the package manager")
	fmt.Println("  list                List installed packages")
	fmt.Println("  help                Show this help message")
}

func doExit(db *sql.DB) {
	err := db.Close()
	if err != nil {
		fmt.Println("Error closing db:", err)
	}
	os.Exit(1)
}

func downloadFile(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func installPackage(packageName string) {
	connStr := "user=postgres password=bezhan2009 dbname=050854c3e9ceb6f18b5978e90b0ff5dcd68517a6d2c48fcc6ece85f63842e6bc host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection to database failed: %v\n", err)
		doExit(db)
	}
	defer db.Close()

	var url string
	query := fmt.Sprintf("SELECT url FROM Packages WHERE name='%s'", packageName)
	err = db.QueryRow(query).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("Package '%s' not found.\n", packageName)
		} else {
			fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
			doExit(db)
		}
		return
	}

	fmt.Printf("Downloading package from: %s\n", url)
	err = downloadFile(url, packageName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		return
	}
	fmt.Printf("Package '%s' installed successfully.\n", packageName)
}

func removePackage(packageName string) {
	filename := "./" + packageName + ".pkg"
	err := os.Remove(filename)
	if err != nil {
		fmt.Printf("Error: Could not remove package '%s'.\n", packageName)
	} else {
		fmt.Printf("Package '%s' removed successfully.\n", packageName)
	}
}

func updatePackageManager() {
	fmt.Println("Updating package manager...")
	updateURL := "your_update_url" // URL для обновления пакетного менеджера
	updateFilename := "opm_new"

	err := downloadFile(updateURL, updateFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		return
	}

	// Замена текущего исполняемого файла новым
	err = os.Rename(updateFilename, "opm")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update package manager: %v\n", err)
	} else {
		fmt.Println("Package manager updated successfully.")
	}
}

func listInstalledPackages() {
	fmt.Println("Listing installed packages:")
	// Для простоты предположим, что установленные пакеты находятся в текущем каталоге с расширением .pkg
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing packages: %v\n", err)
		return
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pkg") {
			fmt.Println(file.Name())
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No command provided.")
		showHelp()
		return
	}

	switch os.Args[1] {
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Error: No package name provided.")
			return
		}
		installPackage(os.Args[2])
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Error: No package name provided.")
			return
		}
		removePackage(os.Args[2])
	case "update":
		updatePackageManager()
	case "list":
		listInstalledPackages()
	case "help":
		showHelp()
	default:
		fmt.Printf("Error: Unknown command '%s'.\n", os.Args[1])
		showHelp()
	}
}
