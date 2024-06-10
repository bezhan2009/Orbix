package ORPXI

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"goCmd/cmdPress"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strings"
)

func CMD(commandInput string) {
	utils.SystemInformation()

	isWorking := true
	isPermission := true

	var promptText string

	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		fmt.Println("Ошибка при проверке директории с паролями:", err)
		return
	}

	if !isEmpty {
		if commandInput == "" {
			dir, _ := os.Getwd()
			user := cmdPress.CmdUser(dir)

			if !CheckUser(user) {
				return
			}
		}
	}

	for isWorking {
		cyan := color.New(color.FgCyan).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		dir, _ := os.Getwd()
		dirC := cmdPress.CmdDir(dir)
		user := cmdPress.CmdUser(dir)

		if promptText != "" {
			fmt.Print("\n" + promptText)
		} else {
			fmt.Printf("\n┌─(%s)-[%s%s]\n", cyan("ORPXI "+user), cyan("~"), cyan(dirC))
			fmt.Printf("└─$ %s", green(">", commandInput))
		}

		commandLine := prompt.Input("", autoComplete)
		commandLine = strings.TrimSpace(commandLine)
		commandParts := parseCommandLine(commandLine)

		if len(commandParts) == 0 {
			continue
		}

		command := commandParts[0]
		commandArgs := commandParts[1:]
		commandLower := strings.ToLower(command)

		commandHistory = append(commandHistory, commandLine)

		fmt.Println()

		if commandLower == "prompt" {
			if len(commandArgs) < 1 {
				fmt.Println("prompt <name_prompt>")
				fmt.Println("to delete prompt enter:")
				fmt.Println("prompt delete")
				continue
			}

			namePrompt := commandArgs[0]

			if namePrompt != "delete" {
				namePrompt = strings.TrimSpace(namePrompt)
				promptText = namePrompt
				fmt.Println("Prompt set to:", promptText)
			} else {
				promptText, _ = os.Getwd()
				fmt.Println("Prompt set to:", promptText)
				promptText = ""
			}

			continue
		}

		helpText := `
Для получения сведений об командах наберите HELP
CREATE             создает новый файл
CLEAN              очистка экрана
CD                 смена текущего каталога
LS                 выводит содержимое каталога
NEWSHABLON         создает новый шаблон комманд для выполнения
REMOVE             удаляет файл
READ               выводит на экран содержимое файла
PROMPT             Изменяет ORPXI.
PINGVIEW           показывает пинг.
PRIMES             Поиск больших простых чисел
PICALC             Вычисление числа π.
NEWUSER            новый пользователь для ORPXI.
ORPXI              запускает ещё одну ORPXI
SHABLON            выполняет определенный шаблон комманд
SYSTEMGOCMD        вывод информации о ORPXI
SYSTEMINFO         вывод информации о системе
SIGNOUT            пользователь выходит из ORPXI
TREE               Графически отображает структуру каталогов диска или пути.
WRITE              записывает данные в файл
EDIT               редактирует файл
WIFIUTILS          Запускает утилиту для работы с WiFi
EXTRACTZIP         распаковывает архивы .zip
SCANPORT           Сканирование портов
WHOIS              Информация о домене
DNSLOOKUP          DNS-запросы
FILEIO             Тест на интенсивную работу с файлами
IPINFO             Информация об IP-адресе
GEOIP              Геолокация IP-адреса
MATRIXMUL          Умножение больших матриц
EXIT               Выход
`

		if commandLower == "help" {
			fmt.Println(helpText)
			continue
		}

		isValid := utils.ValidCommand(commandLower, commands)

		if !isValid {
			fullCommand := append([]string{command}, commandArgs...)
			err := utils.ExternalCommand(fullCommand)
			if err != nil {
				fullPath := filepath.Join(dir, command)
				fullCommand[0] = fullPath
				err = utils.ExternalCommand(fullCommand)
				if err != nil {
					suggestedCommand := suggestCommand(commandLower)
					fmt.Printf("Ошибка при запуске команды '%s': %v\n", commandLine, err)
					if suggestedCommand != "" {
						fmt.Printf("Возможно, вы имели в виду: %s?\n", suggestedCommand)
					}
				}
			}
			continue
		}

		if commandInput != "" {
			isPermission = false
		} else {
			isPermission = true
		}

		ExecuteCommand(commandLower, command, commandLine, dir, commands, commandArgs, &isWorking, isPermission)
	}
}
