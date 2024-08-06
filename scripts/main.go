package main

import (
	"fmt"
	"goCmd/run"
	"os"
	"os/exec"
	"strings"
)

// checkSystemResources checks the system resources.
func checkSystemResources() error {
	// Check CPU usage and memory
	cpuUsage, err := getCPUUsage()
	if err != nil {
		return fmt.Errorf("error getting CPU usage: %v", err)
	}
	if cpuUsage > 80 {
		return fmt.Errorf("high CPU usage: %f%%", cpuUsage)
	}

	//memStats := runtime.MemStats{}
	//runtime.ReadMemStats(&memStats)
	//if memStats.Frees < 100*1024*1024 { // less than 100 MB free memory
	//	return fmt.Errorf("low free memory: %d bytes", memStats.Frees)
	//}

	// Check free disk space
	//freeDiskSpace, err := getFreeDiskSpace("C:")
	//if err != nil {
	//	return fmt.Errorf("error getting free disk space: %v", err)
	//}
	//if freeDiskSpace < 10*1024*1024*1024 { // less than 10 GB free space
	//	return fmt.Errorf("low free disk space: %d bytes", freeDiskSpace)
	//}

	return nil
}

// getCPUUsage returns the current CPU usage in percentage.
func getCPUUsage() (float64, error) {
	// Depends on the OS; for Windows, use PowerShell or WMI.
	cmd := exec.Command("wmic", "cpu", "get", "loadpercentage")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("unable to get CPU load data")
	}
	cpuUsage := strings.TrimSpace(lines[1])
	var usage float64
	_, err = fmt.Sscanf(cpuUsage, "%f", &usage)
	if err != nil {
		return 0, err
	}
	return usage, nil
}

// getFreeDiskSpace returns the free disk space in bytes.
func getFreeDiskSpace(drive string) (uint64, error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("(Get-PSDrive -Name %s).Used", drive))
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	var usedSpace uint64
	_, err = fmt.Sscanf(string(output), "%d", &usedSpace)
	if err != nil {
		return 0, err
	}
	return usedSpace, nil
}

// checkForBlockingProcesses checks for processes that may interfere.
func checkForBlockingProcesses() error {
	blockingProcesses := []string{"example.exe"} // List of processes that may interfere
	for _, process := range blockingProcesses {
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", process))
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("error checking processes: %v", err)
		}
		if strings.Contains(string(output), process) {
			return fmt.Errorf("interfering process found: %s", process)
		}
	}
	return nil
}

func main() {
	// Check system resources
	if err := checkSystemResources(); err != nil {
		fmt.Println("Error checking system resources:", err)
		os.Exit(1) // Exit with error code
	}

	// Check for interfering processes
	if err := checkForBlockingProcesses(); err != nil {
		fmt.Println("Error checking processes:", err)
		os.Exit(1) // Exit with error code
	}

	// Start the main program
	run.Init()
}
