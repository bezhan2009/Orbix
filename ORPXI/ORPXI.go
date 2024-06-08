package ORPXI

import (
	"bufio"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	"goCmd/Network"
	"goCmd/Network/wifiUtils"
	"goCmd/cmdPress"
	"goCmd/commands/commandsWithSignaiture/Create"
	"goCmd/commands/commandsWithSignaiture/Edit"
	"goCmd/commands/commandsWithSignaiture/ExtractZip"
	"goCmd/commands/commandsWithSignaiture/Read"
	"goCmd/commands/commandsWithSignaiture/Remove"
	"goCmd/commands/commandsWithSignaiture/Rename"
	"goCmd/commands/commandsWithSignaiture/Write"
	"goCmd/commands/commandsWithSignaiture/shablon"
	"goCmd/commands/commandsWithoutSignature/CD"
	"goCmd/commands/commandsWithoutSignature/Clean"
	"goCmd/commands/commandsWithoutSignature/Ls"
	"goCmd/commands/resourceIntensive/FileIOStressTest"
	"goCmd/commands/resourceIntensive/MatrixMultiplication"
	"goCmd/commands/resourceIntensive/PiCalculation"
	"goCmd/commands/resourceIntensive/PrimeNumbers"
	"goCmd/debug"
	"goCmd/structs"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var commands = []structs.Command{
	{"whois", "Информация о домене"},
	{"pingview", "Показывает пинг"},
	{"traceroute", "Трассировка маршрута"},
	{"extractzip", "Распаковывает архивы .zip"},
	{"signout", "Пользователь выходит из ORPXI"},
	{"newshablon", "Создает новый шаблон комманд для выполнения"},
	{"shablon", "Выполняет определенный шаблон комманд"},
	{"newuser", "Новый пользователь для ORPXI"},
	{"promptSet", "Изменяет ORPXI"},
	{"systemgocmd", "Вывод информации о ORPXI"},
	{"rename", "Переименовывает файл"},
	{"remove", "Удаляет файл"},
	{"read", "Выводит на экран содержимое файла"},
	{"write", "Записывает данные в файл"},
	{"create", "Создает новый файл"},
	{"exit", "Выход из программы"},
	{"orpxi", "Запускает ещё одну ORPXI"},
	{"wifiutils", "Запускает утилиту для работы с WiFi"},
	{"clean", "Очистка экрана"},
	{"matrixmul", "Умножение больших матриц"},
	{"primes", "Поиск больших простых чисел"},
	{"picalc", "Вычисление числа π"},
	{"fileio", "Тест на интенсивную работу с файлами"},
	{"cd", "Смена текущего каталога"},
	{"edit", "Редактирует файл"},
	{"ls", "Выводит содержимое каталога"},
	{"scanport", "Сканирование портов"},
	{"dnslookup", "DNS-запросы"},
	{"ipinfo", "Информация об IP-адресе"},
	{"geoip", "Геолокация IP-адреса"},
}

var commandHistory []string

func autoComplete(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	if len(text) == 0 {
		return []prompt.Suggest{}
	}

	parts := strings.Split(text, " ")
	if len(parts) == 1 {
		// Suggest command names
		return prompt.FilterHasPrefix(createUniqueCommandSuggestions(), text, true)
	} else {
		// Suggest file or directory names
		dir := "."
		if len(parts) > 2 {
			dir = strings.Join(parts[:len(parts)-1], " ")
		}
		return prompt.FilterHasPrefix(createFileSuggestions(dir), parts[len(parts)-1], true)
	}
}

func createUniqueCommandSuggestions() []prompt.Suggest {
	uniqueCommands := make(map[string]struct{})
	var suggestions []prompt.Suggest

	for _, cmd := range commands {
		if _, exists := uniqueCommands[cmd.Name]; !exists {
			uniqueCommands[cmd.Name] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: cmd.Name, Description: cmd.Description})
		}
	}
	for _, cmd := range commandHistory {
		if _, exists := uniqueCommands[cmd]; !exists {
			uniqueCommands[cmd] = struct{}{}
			suggestions = append(suggestions, prompt.Suggest{Text: cmd})
		}
	}

	return suggestions
}

func createFileSuggestions(dir string) []prompt.Suggest {
	files, err := os.ReadDir(dir)
	if err != nil {
		return []prompt.Suggest{}
	}

	var suggestions []prompt.Suggest
	for _, file := range files {
		suggestions = append(suggestions, prompt.Suggest{Text: file.Name()})
	}
	return suggestions
}

func CMD(commandInput string) {
	utils.SystemInformation()

	isWorking := true
	isPermission := true

	var promptText string

	// Проверка пароля
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
			fmt.Printf("└─$ %s", green("> "))
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

		executeCommand(commandLower, command, commandLine, dir, commands, commandArgs, &isWorking, isPermission)
	}
}

func parseCommandLine(commandLine string) []string {
	var parts []string
	var currentPart strings.Builder
	var inQuotes bool

	for _, char := range commandLine {
		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				currentPart.WriteRune(char)
			} else {
				parts = append(parts, currentPart.String())
				currentPart.Reset()
			}
		default:
			currentPart.WriteRune(char)
		}
	}

	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	return parts
}

func suggestCommand(input string) string {
	for _, cmd := range commands {
		if strings.HasPrefix(cmd.Name, input) {
			return cmd.Name
		}
	}
	return ""
}

func executeCommand(commandLower string, command string, commandLine string, dir string, commands []structs.Command, commandArgs []string, isWorking *bool, isPermission bool) {
	user := cmdPress.CmdUser(dir)
	switch commandLower {
	case "wifiutils":
		wifiUtils.Start()
	case "pingview":
		Network.Ping(commandArgs)
	case "traceroute":
		Network.Traceroute(commandArgs)
	case "extractzip":
		if len(commandArgs) < 2 {
			fmt.Println("Usage: extractzip <zipfile> <destination>")
		} else {
			err := ExtractZip.ExtractZip(commandArgs[0], commandArgs[1])
			if err != nil {
				fmt.Println("Error extracting ZIP file:", err)
			}
		}
	case "scanport":
		if len(commandArgs) < 2 {
			fmt.Println("Usage: scanport <host> <ports>")
		} else {
			ports := []int{}
			for _, p := range commandArgs[1:] {
				port, err := strconv.Atoi(p)
				if err != nil {
					fmt.Printf("Invalid port: %s\n", p)
					return
				}
				ports = append(ports, port)
			}
			Network.ScanPort(commandArgs[0], ports)
		}
	case "whois":
		if len(commandArgs) < 1 {
			fmt.Println("Usage: whois <domain>")
		} else {
			Network.Whois(commandArgs[0])
		}
	case "dnslookup":
		if len(commandArgs) < 1 {
			fmt.Println("Usage: dnslookup <domain>")
		} else {
			Network.DNSLookup(commandArgs[0])
		}
	case "ipinfo":
		if len(commandArgs) < 1 {
			fmt.Println("Usage: ipinfo <ip>")
		} else {
			Network.IPInfo(commandArgs[0])
		}
	case "geoip":
		if len(commandArgs) < 1 {
			fmt.Println("Usage: geoip <ip>")
		} else {
			Network.GeoIP(commandArgs[0])
		}
	case "orpxi":
		if isPermission {
			CMD("")
		}
	case "newuser":
		if isPermission {
			NewUser()
		}
	case "signout":
		if isPermission {
			if !CheckUser(user) {
				*isWorking = false
			}
		}
	case "matrixmul":
		MatrixMultiplication.MatrixMulCommand()
	case "primes":
		limit := 100000
		if len(commandArgs) > 0 {
			l, err := strconv.Atoi(commandArgs[0])
			if err == nil {
				limit = l
			}
		}
		PrimeNumbers.PrimeCommand(limit)
	case "picalc":
		precision := 10000
		if len(commandArgs) > 0 {
			p, err := strconv.Atoi(commandArgs[0])
			if err == nil {
				precision = p
			}
		}
		PiCalculation.PiCalcCommand(precision)
	case "fileio":
		filename := "largefile.dat"
		size := 100 * 1024 * 1024
		if len(commandArgs) > 0 {
			s, err := strconv.Atoi(commandArgs[0])
			if err == nil {
				size = s
			}
		}
		FileIOStressTest.FileIOCommand(filename, size)
	case "newshablon":
		shablon.Make()
	case "shablon":
		if len(commandArgs) < 1 {
			fmt.Println("Использование: shablon <название_шаблона>")
			return
		}

		nameShablon := commandArgs[0]
		err := Start(nameShablon)
		if err != nil {
			fmt.Println(err)
		}
	case "systemgocmd":
		utils.SystemInformation()
	case "exit":
		if isPermission {
			*isWorking = false
		}
	case "create":
		name, err := Create.File(commandArgs)
		if err != nil {
			fmt.Println(err)
			debug.Commands(command, false)
		} else if name != "" {
			fmt.Printf("Файл %s успешно создан!\n", name)
			fmt.Printf("Директория нового файла: %s\n", filepath.Join(dir, name))
			debug.Commands(command, true)
		}
	case "write":
		Write.File(commandLower, commandArgs)
	case "read":
		Read.File(commandLower, commandArgs)
	case "remove":
		name, err := Remove.File(commandArgs)
		if err != nil {
			debug.Commands(command, false)
			fmt.Println(err)
		} else {
			debug.Commands(command, true)
			fmt.Printf("Файл %s успешно удален!\n", name)
		}
	case "rename":
		errRename := Rename.Rename(commandArgs)
		if errRename != nil {
			debug.Commands(command, false)
			fmt.Println(errRename)
		} else {
			debug.Commands(command, true)
		}
	case "clean":
		Clean.Screen()
	case "cd":
		if len(commandArgs) == 0 {
			dir, _ := os.Getwd()
			fmt.Println(dir)
		} else {
			err := CD.ChangeDirectory(commandArgs[0])
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	case "edit":
		if len(commandArgs) < 1 {
			fmt.Println("Использование: edit <файл>")
			return
		}
		filename := commandArgs[0]
		err := Edit.File(filename)
		if err != nil {
			fmt.Println(err)
		}
	case "ls":
		Ls.PrintLS()
	default:
		validCommand := utils.ValidCommand(commandLower, commands)
		if !validCommand {
			fmt.Printf("'%s' не является внутренней или внешней командой,\nисполняемой программой или пакетным файлом.\n", commandLine)
			suggestedCommand := suggestCommand(commandLower)
			if suggestedCommand != "" {
				fmt.Printf("Возможно, вы имели в виду: %s?\n", suggestedCommand)
			}
		}
	}
}

func Start(shablonName string) error {
	shablonName = strings.TrimSpace(shablonName)

	file, err := os.OpenFile(shablonName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Игнорировать пустые строки
		}
		CMD(line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}
