package src

import (
	"fmt"
	"goCmd/cmd/dirInfo"
	"goCmd/src/environment"
	"goCmd/system"
	"os"
	"strings"
	"time"
)

func printPromptInfo(location, user, dirC, commandInput *string, sd *system.Session) {
	// Обрезаем Prompt, если он длинный
	if len(system.Prompt) > 2 {
		system.Prompt = string(system.Prompt[:2])
	}

	// Сохраняем форматированные данные в переменные
	gitBranch := system.Green(sd.GitBranch)
	userInfo := system.Cyan("Orbix@" + getUser(*user))
	locationInfo := system.Yellow(*location)
	dirInfo := system.Cyan(*dirC)
	currentTime := system.Magenta(time.Now().Format("15:04"))
	prompt := strings.TrimSpace(system.Prompt)
	input := strings.TrimSpace(*commandInput)

	// Формируем строки для вывода
	header := fmt.Sprintf(
		"\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s%s%s%s%s",
		system.Yellow("╭"), system.Yellow("─"), system.Yellow("("),
		userInfo, system.Yellow(")"), system.Yellow("─"), system.Yellow("["),
		locationInfo, currentTime, system.Yellow("]"), system.Yellow("─"), system.Yellow("["),
		system.Cyan("~"), dirInfo, system.Yellow("]"), system.Yellow(" git:"), system.Yellow("["), gitBranch, system.Yellow("]"),
	)

	footer := fmt.Sprintf("%s%s %s", system.Yellow("╰"), system.Green(prompt), system.Green(*commandInput))

	// Печатаем информацию
	fmt.Println(header)
	fmt.Print(footer)

	// Выводим перенос строки, если есть команды
	if input != "" && len(os.Args) > 0 {
		fmt.Println()
	}
}

func PrintPromptInfoWithoutGit(location, user, dirC, commandInput *string) {
	// Обрезаем Prompt, если он длинный
	if len(system.Prompt) > 2 {
		system.Prompt = string(system.Prompt[:1])
	}

	// Сохраняем форматированные данные в переменные
	userInfo := system.Cyan("Orbix@" + getUser(*user))
	locationInfo := system.Yellow(*location)
	dirInfo := system.Cyan(*dirC)
	currentTime := system.Magenta(time.Now().Format("15:04"))
	prompt := strings.TrimSpace(system.Prompt)

	// Формируем строки для вывода
	header := fmt.Sprintf(
		"\n%s%s%s%s%s%s%s%s %s%s%s%s%s%s%s",
		system.Yellow("╭"), system.Yellow("─"), system.Yellow("("),
		userInfo, system.Yellow(")"), system.Yellow("─"), system.Yellow("["),
		locationInfo, currentTime, system.Yellow("]"), system.Yellow("─"), system.Yellow("["),
		system.Cyan("~"), dirInfo, system.Yellow("]"),
	)

	footer := fmt.Sprintf("%s%s %s", system.Yellow("╰"), system.Green(prompt), system.Green(*commandInput))

	// Печатаем информацию
	fmt.Println(header)
	fmt.Print(footer)

	// Выводим перенос строки, если есть команды
	if *commandInput != "" && len(os.Args) > 0 {
		fmt.Println()
	}
}

func customPrompt(commandInput, prompt *string,
	colorsMap map[string]func(...interface{}) string) {
	if strings.TrimSpace(*commandInput) != "" {
		splitPrompt := strings.Split(*prompt, ", ")
		fmt.Printf("%s%s", colorsMap[splitPrompt[1]](splitPrompt[0]), system.Green(*commandInput))
	} else {
		splitPrompt := strings.Split(*prompt, ", ")
		fmt.Print(colorsMap[splitPrompt[1]](splitPrompt[0]))
	}
}

func printOldPrompt(commandInput, dir *string) {
	if strings.TrimSpace(*commandInput) != "" {
		fmt.Printf("ORB %s> %s", *dir, system.Green(*commandInput))
	} else {
		fmt.Printf("ORB %s> ", *dir)
	}
}

func OrbixPrompt(session *system.Session,
	prompt, commandInput *string,
	isWorking, isPermission *bool,
	colorsMap *map[string]func(...interface{}) string) {
	if session.IsAdmin {
		if *prompt == "" {
			printOldPrompt(commandInput, &system.UserDir)
		} else {
			customPrompt(commandInput, prompt,
				*colorsMap)
		}

		return
	}

	Orbixuser, _ := environment.GetVariableValue("user")
	if Orbixuser == "" {
		Orbixuser = dirInfo.CmdUser(&system.UserDir)
	}

	OrbixuserStr := fmt.Sprintf("%s", Orbixuser)

	if !session.IsAdmin {
		dirC = dirInfo.CmdDir(system.UserDir)

		// Single user check outside repeated prompt formatting
		if !system.Unauthorized {
			go func() {
				watchFile(system.RunningPath, OrbixuserStr, isWorking, isPermission)
			}()
		}

		if *prompt == "" {
			if system.GitExists {
				printPromptInfo(&system.Location,
					&OrbixuserStr,
					&dirC,
					commandInput,
					session) // New helper function
			} else {
				PrintPromptInfoWithoutGit(&system.Location,
					&OrbixuserStr,
					&dirC,
					commandInput) // New helper function
			}
		} else {
			customPrompt(commandInput, prompt,
				*colorsMap)
		}
	}
}
