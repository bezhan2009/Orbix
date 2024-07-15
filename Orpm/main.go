package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lopm
#include "opm.h"
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func showHelp() {
	C.show_help()
}

func installPackage(packageName string) {
	cPackageName := C.CString(packageName)
	defer C.free(unsafe.Pointer(cPackageName))
	C.install_package(cPackageName)
}

func removePackage(packageName string) {
	cPackageName := C.CString(packageName)
	defer C.free(unsafe.Pointer(cPackageName))
	C.remove_package(cPackageName)
}

func updatePackageManager() {
	C.update_package_manager()
}

func listInstalledPackages() {
	C.list_installed_packages()
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
