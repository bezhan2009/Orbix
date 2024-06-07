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
	"time"
)

func animatedPrint(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(2 * time.Millisecond)
	}
}

func CMD(commandInput string) {
	utils.SystemInformation()

	isWorking := true
	isPermission := true
	reader := bufio.NewReader(os.Stdin)

	var prompt string

	// Проверка пароля
	isEmpty, err := isPasswordDirectoryEmpty()
	if err != nil {
		animatedPrint("Ошибка при проверке директории с паролями:" + err.Error() + "\n")
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
			animatedPrint("\n" + prompt)
		} else {
			animatedPrint(fmt.Sprintf("\n┌─(%s)-[%s%s]\n", cyan("ORPXI "+user), cyan("~"), cyan(dirC)))
			animatedPrint(fmt.Sprintf("└─$ %s", green(commandInput)))
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

		animatedPrint("\n")

		if commandLower == "prompt" {
			if len(commandArgs) < 1 {
				animatedPrint("prompt <name_prompt>\n")
				animatedPrint("to delete prompt enter:\n")
				animatedPrint("prompt delete\n")
				continue
			}

			namePrompt := commandArgs[0]

			if namePrompt != "delete" {
				namePrompt = strings.TrimSpace(namePrompt)
				prompt = namePrompt
				animatedPrint(fmt.Sprintf("Prompt set to: %s\n", prompt))
			} else {
				prompt, _ = os.Getwd()
				animatedPrint(fmt.Sprintf("Prompt set to: %s\n", prompt))
				prompt = ""
			}

			continue
		}

		if commandLower == "help" {
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
NEWUSER            новый пользователь для ORPXI.
ORPXI              запускает ещё одну ORPXI
SHABLON            выполняет определенный шаблон комманд
SYSTEMGOCMD        вывод информации о ORPXI
SYSTEMINFO         вывод информации о системе
SIGNOUT            пользователь выходит из ORPXI
TREE               Графически отображает структуру каталогов диска или пути.
WRITE              записывает данные в файл
EDIT               редактирует файл
EXTRACTZIP         распаковывает архивы .zip
EXIT               Выход
`
			animatedPrint(helpText)

			errDebug := debug.Commands(command, true)
			if errDebug != nil {
				animatedPrint(errDebug.Error() + "\n")
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
					animatedPrint(fmt.Sprintf("Ошибка при запуске команды '%s': %v\n", commandLine, err))
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
			animatedPrint("Usage: extractzip <zipfile> <destination>\n")
		} else {
			err := ExtractZip.ExtractZip(commandArgs[0], commandArgs[1])
			if err != nil {
				animatedPrint("Error extracting ZIP file: " + err.Error() + "\n")
			}
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
	case "newshablon":
		shablon.Make()
	case "shablon":
		if len(commandArgs) < 1 {
			animatedPrint("Использования: shablon <названия_шаблона>\n")
			return
		}

		nameShablon := commandArgs[0]
		err := Start(nameShablon)
		if err != nil {
			animatedPrint(err.Error() + "\n")
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
			animatedPrint(err.Error() + "\n")
			debug.Commands(command, false)
		} else if name != "" {
			animatedPrint(fmt.Sprintf("Файл %s успешно создан!!!\n", name))
			animatedPrint(fmt.Sprintf("Директория нового файла: %s\n", filepath.Join(dir, name)))
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
			animatedPrint(err.Error() + "\n")
		} else {
			debug.Commands(command, true)
			animatedPrint(fmt.Sprintf("Файл %s успешно удален!!!\n", name))
		}
	case "rename":
		errRename := Rename.Rename(commandArgs)
		if errRename != nil {
			debug.Commands(command, false)
			animatedPrint(errRename.Error() + "\n")
		} else {
			debug.Commands(command, true)
		}
	case "clean":
		Clean.Screen()
	case "cd":
		if len(commandArgs) == 0 {
			dir, _ := os.Getwd()
			animatedPrint(dir + "\n")
		} else {
			err := CD.ChangeDirectory(commandArgs[0])
			if err != nil {
				animatedPrint(err.Error() + "\n")
			}
			return
		}
	case "edit":
		if len(commandArgs) < 1 {
			animatedPrint("Использование: edit <файл>\n")
			return
		}
		filename := commandArgs[0]
		err := Edit.File(filename)
		if err != nil {
			animatedPrint(err.Error() + "\n")
		}
	case "ls":
		Ls.PrintLS()
	default:
		validCommand := utils.ValidCommand(commandLower, commands)
		if !validCommand {
			animatedPrint(fmt.Sprintf("'%s' не является внутренней или внешней командой,\nисполняемой программой или пакетным файлом.\n", commandLine))
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
