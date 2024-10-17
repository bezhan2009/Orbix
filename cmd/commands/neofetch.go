package commands

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"goCmd/system"
	"os"
	"runtime"
	"time"
)

func FetchNeofetch(user string) {
	// Получаем информацию о пользователе
	username := user

	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()

	// Получаем информацию о системе
	osName := fmt.Sprintf("%s %s", system.OperationSystem, runtime.GOARCH)
	hostStat, _ := host.Info()
	uptime := fmt.Sprintf("%v", time.Duration(hostStat.Uptime)*time.Second)
	cpuInfo, _ := cpu.Info()
	cpuModel := cpuInfo[0].ModelName
	memStat, _ := mem.VirtualMemory()
	memory := fmt.Sprintf("%.2fMiB / %.2fMiB", float64(memStat.Used)/1024/1024, float64(memStat.Total)/1024/1024)

	// Проверяем переменные окружения и устанавливаем значения по умолчанию
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash" // Или другой шела по умолчанию
	}
	terminal := os.Getenv("TERM")
	if terminal == "" {
		terminal = "unknown" // Или определить другой терминал по умолчанию
	}

	// Цветное оформление
	title := color.New(color.FgCyan, color.Bold).SprintFunc()
	info := color.New(color.FgWhite).SprintFunc()
	value := color.New(color.FgYellow).SprintFunc()

	// ASCII-иконка Windows
	windows_10 := []string{
		"                                ..,",
		"                    ....,,:;+ccllll",
		"      ...,,+:;  cllllllllllllllllll",
		",cclllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"                                   ",
		"llllllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"llllllllllllll  lllllllllllllllllll",
		"`'ccllllllllll  lllllllllllllllllll",
		"       `' \\*::  :ccllllllllllllllll",
		"                       ````''*::cll",
		"                                 ``",
	}

	neofetch := fmt.Sprintf(`%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s`,
		cyan(fmt.Sprintf("       %s           %s@%s", title(windows_10[0]), title(username), title("user"))),
		cyan(fmt.Sprintf("       %s           %s", title(windows_10[1]), info("-----------"))),
		cyan(fmt.Sprintf("       %s           OS: %s", title(windows_10[2]), value(osName))),
		cyan(fmt.Sprintf("       %s           Kernel: %s", title(windows_10[3]), value(hostStat.KernelVersion))),
		cyan(fmt.Sprintf("       %s           Terminal: %s", title(windows_10[4]), value(terminal))),
		cyan(fmt.Sprintf("       %s           Shell: %s", title(windows_10[5]), value(shell))),
		cyan(fmt.Sprintf("       %s           Uptime: %s", title(windows_10[6]), value(uptime))),
		cyan(fmt.Sprintf("       %s           CPU: %s", title(windows_10[7]), value(cpuModel))),
		cyan(fmt.Sprintf("       %s           Memory: %s", title(windows_10[8]), value(memory))),
		cyan(fmt.Sprintf("       %s           ", title(windows_10[9]))),
		cyan(fmt.Sprintf("       %s           ", title(windows_10[10]))),
		cyan(fmt.Sprintf("       %s   ", title(windows_10[11]))),
		cyan(fmt.Sprintf("       %s   ", title(windows_10[12]))),
		cyan(fmt.Sprintf("       %s   ", title(windows_10[13]))),
		cyan(fmt.Sprintf("       %s   ", title(windows_10[14]))),
		cyan(fmt.Sprintf("       %s   ", title(windows_10[15]))),
		cyan(fmt.Sprintf("       %s   ", title(windows_10[16]))),
		cyan(fmt.Sprintf("       %s   ", title(windows_10[17]))),
		cyan(fmt.Sprintf("       %s   ", title(windows_10[18]))),
		cyan(fmt.Sprintf("       %s   ", title(windows_10[19]))))

	fmt.Println(neofetch)
	fmt.Println()
}
