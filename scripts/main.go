package main

import (
	"fmt"
	"goCmd/run"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// checkSystemResources checks the system resources.
func checkSystemResources() error {
	// Check CPU usage
	cpuUsage, err := getCPUUsage()
	if err != nil {
		return fmt.Errorf("error getting CPU usage: %v", err)
	}
	if cpuUsage > 80 {
		return fmt.Errorf("high CPU usage: %f%%", cpuUsage)
	}

	// You can add memory and disk space checks similarly
	// ...

	return nil
}

// getCPUUsage returns the current CPU usage in percentage.
func getCPUUsage() (float64, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("wmic", "cpu", "get", "loadpercentage")
	case "linux":
		cmd = exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)' | sed 's/.*, *\\([0-9.]*\\)%* id.*/\\1/' | awk '{print 100 - $1}'")
	case "darwin":
		cmd = exec.Command("sh", "-c", "ps -A -o %cpu | awk '{s+=$1} END {print s}'")
	default:
		return 0, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	return parseCPUUsageOutput(output)
}

// parseCPUUsageOutput parses the output from the CPU usage command.
func parseCPUUsageOutput(output []byte) (float64, error) {
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 1 {
		return 0, fmt.Errorf("unable to get CPU load data")
	}
	// For Windows
	if runtime.GOOS == "windows" {
		if len(lines) < 2 {
			return 0, fmt.Errorf("unable to get CPU load data")
		}
		return strconv.ParseFloat(strings.TrimSpace(lines[1]), 64)
	}
	// For Linux and macOS
	cpuUsage := strings.TrimSpace(lines[len(lines)-1])
	return strconv.ParseFloat(cpuUsage, 64)
}

// getFreeDiskSpace returns the free disk space in bytes.
func getFreeDiskSpace(drive string) (uint64, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("(Get-PSDrive -Name %s).Free", drive))
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", fmt.Sprintf("df -k %s | tail -1 | awk '{print $4}'", drive))
	default:
		return 0, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	var freeSpace uint64
	_, err = fmt.Sscanf(string(output), "%d", &freeSpace)
	if err != nil {
		return 0, err
	}
	return freeSpace * 1024, nil // Convert from KB to bytes
}

// checkForBlockingProcesses checks for processes that may interfere.
func checkForBlockingProcesses() error {
	blockingProcesses := []string{"example.exe"} // List of processes that may interfere
	for _, process := range blockingProcesses {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", process))
		case "linux", "darwin":
			cmd = exec.Command("sh", "-c", fmt.Sprintf("ps -A | grep %s", process))
		default:
			return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
		}
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
