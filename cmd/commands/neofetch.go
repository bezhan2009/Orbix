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

func FetchNeofetch(GlobalSession *system.Session) {
	// Получаем информацию о пользователе
	username := GlobalSession.User

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
%s`,
		cyan(fmt.Sprintf("       /0d/+00sssooyPONyssss/   /0-/+oosss0oyMMNyssssD           %s@%s", title(username), title("user"))),
		cyan(fmt.Sprintf("       /ef/+oossso00ygSNySss/   /0-/+F0sssooyMMNyssssG           %s", info("-----------"))),
		cyan(fmt.Sprintf("       /ef/+oossso00ygSNySss/   /0-/+F0sssooyMMNyssssG           OS: %s", value(osName))),
		cyan(fmt.Sprintf("       /ef/+oossso00ygSNySss/   /0-/+F0sssooyMMNyssssG           Kernel: %s", value(hostStat.KernelVersion))),
		cyan(fmt.Sprintf("       /9-/+oosss0yGGNyssss-+   /0-/+oo00sooyMMNyssss+           Terminal: %s", value(terminal))),
		cyan(fmt.Sprintf("       /P=/+00sss0oySFNyssss/   /0-/+D0sssooyMMNyssss\\           Shell: %s", value(shell))),
		cyan(fmt.Sprintf("       /6=/+oosso0oyADNyssgs/   /0-/+oSOssoayDFNyssss-           Uptime: %s", value(uptime))),
		cyan(fmt.Sprintf("       /+Y/+oosssooyLDNyssss/   /0-/+ooFssooyMMNyssss/           CPU: %s", value(cpuModel))),
		cyan(fmt.Sprintf("                                                                 Memory: %s", value(memory))),
		cyan(fmt.Sprintf("       /0-/+ooss=aoyMMNyssasa   /0-/+oossso+yOOOyssss/           ")),
		cyan(fmt.Sprintf("       /6=/+oosso0oyADNyssgs/   /0-/+oSOssoayDFNyssss-           ")),
		cyan("       /0-/+oosssooyMMNyssa-+   /0-/+oosss(0yMMNyssss/   "),
		cyan("       /0-/+ooss--oyMMNysfss/   /0-/+odsssfsyMMNyssss   "),
		cyan("       /0-/+oosssooyMMNydsss+   /0-/+ofsss+-yMMNyssss/   "),
		cyan("       /0-/+oosssooyMMNyssa-+   /0-/+oosss(0yMMNyssss/   "),
		cyan("       /0-/+oosssooyMMNyssa-+   /0-/+oosss(0yMMNyssss/   "),
		cyan("       /0-/+oosssooyMMNyssds=   /0-/+oosssooyMMNyssss/"))

	fmt.Println(neofetch)
}
