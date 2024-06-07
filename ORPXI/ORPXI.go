package ORPXI

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"goCmd/Network"
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
	"goCmd/debug"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strings"
)

func CMD(commandInput string) {
	utils.SystemInformation()

	isWorking := true
	isPermission := true
	reader := bufio.NewReader(os.Stdin)

	var prompt string

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

		if prompt != "" {
			fmt.Printf("\n%s", prompt)
		} else {
			fmt.Printf("\n┌─(%s)-[%s%s]\n", cyan("ORPXI "+user), cyan("~"), cyan(dirC))
			fmt.Printf("└─$ %s", green(commandInput))
		}
		var commandLine string
		var commandParts []string
		var commandArgs []string
		var commandLower string
		var command string
		if commandInput != "" {
			isWorking = false
			isPermission = false
			commandLine = strings.TrimSpace(commandInput) // Исправлено использование commandInput
			commandParts = utils.SplitCommandLine(commandLine)
			if len(commandParts) == 0 {
				continue
			}

			command = commandParts[0]
			commandArgs = commandParts[1:]
			commandLower = strings.ToLower(command)
		} else {
			commandLine, _ = reader.ReadString('\n')
			commandLine = strings.TrimSpace(commandLine)
			commandParts = utils.SplitCommandLine(commandLine)

			if len(commandParts) == 0 {
				continue
			}

			command = commandParts[0]
			commandArgs = commandParts[1:]
			commandLower = strings.ToLower(command)
		}

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
				prompt = namePrompt
				fmt.Printf("Prompt set to: %s\n", prompt)
			} else {
				prompt, _ = os.Getwd()
				fmt.Printf("Prompt set to: %s\n", prompt)
				prompt = ""
			}

			continue
		}

		if commandLower == "help" {
			fmt.Println("Для получения сведений об командах наберите HELP")
			fmt.Println("CREATE             создает новый файл")
			fmt.Println("CLEAN              очистка экрана")
			fmt.Println("CD                 смена текущего каталога")
			fmt.Println("LS                 выводит содержимое каталога")
			fmt.Println("NEWSHABLON         создает новый шаблон комманд для выполнения")
			fmt.Println("REMOVE             удаляет файл")
			fmt.Println("READ               выводит на экран содержимое файла")
			fmt.Println("PROMPT             Изменяет ORPXI.")
			fmt.Println("PINGVIEW           показывает пинг.")
			fmt.Println("NEWUSER            новый пользователь для ORPXI.")
			fmt.Println("ORPXI              запускает ещё одну ORPXI")
			fmt.Println("SHABLON            выполняет определенный шаблон комманд")
			fmt.Println("SYSTEMGOCMD        вывод информации о ORPXI")
			fmt.Println("SYSTEMINFO         вывод информации о системе")
			fmt.Println("SIGNOUT            пользователь выходит из ORPXI")
			fmt.Println("TREE               Графически отображает структуру каталогов диска или пути.")
			fmt.Println("WRITE              записывает данные в файл")
			fmt.Println("EDIT               редактирует файл")
			fmt.Println("EXTRACTZIP         распаковывает архивы .zip")
			fmt.Println("EXIT               Выход")

			errDebug := debug.Commands(command, true)
			if errDebug != nil {
				fmt.Println(errDebug)
			}
			continue
		}

		if commandLower == "help" {
			continue
		}

		commands := []string{"pingview", "tracerout", "extractzip", "signout", "newshablon", "shablon", "newuser", "promptSet", "systemgocmd", "rename", "remove", "read", "write", "create", "exit", "orpxi", "clean", "cd", "edit", "ls"}

		isValid := utils.ValidCommand(commandLower, commands)

		if !isValid {
			fullCommand := append([]string{command}, commandArgs...)
			err := utils.ExternalCommand(fullCommand)
			if commandLower == "help" {
				continue
			}
			if err != nil {
				fullPath := filepath.Join(dir, command)
				fullCommand[0] = fullPath
				err = utils.ExternalCommand(fullCommand)
				if err != nil {
					fmt.Printf("Ошибка при запуске команды '%s': %v\n", commandLine, err)
				}
			}
			continue
		}

		executeCommand(commandLower, command, commandLine, dir, commands, commandArgs, &isWorking, isPermission)

	}
}

func executeCommand(commandLower string, command string, commandLine string, dir string, commands []string, commandArgs []string, isWorking *bool, isPermission bool) {
	user := cmdPress.CmdUser(dir)
	switch commandLower {
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
	case "orpxi":
		if isPermission {
			CMD("")
		}
	case "password":
		if isPermission {
			Password()
		}
	case "signout":
		if isPermission {
			if !CheckUser(user) {
				*isWorking = false
			}
		}
	case "newshablon":
		shablon.Make()
	case "shablon":
		if len(commandArgs) < 1 {
			fmt.Println("Использования: shablon <названия_шаблона>")
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
			fmt.Printf("Файл %s успешно создан!!!\n", name)
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
			fmt.Printf("Файл %s успешно удален!!!\n", name)
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
