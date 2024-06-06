package ORPXI

import (
	"bufio"
	"fmt"
	"goCmd/commands/commandsWithSignaiture/Create"
	"goCmd/commands/commandsWithSignaiture/Edit"
	"goCmd/commands/commandsWithSignaiture/Read"
	"goCmd/commands/commandsWithSignaiture/Remove"
	"goCmd/commands/commandsWithSignaiture/Rename"
	"goCmd/commands/commandsWithSignaiture/Write"
	"goCmd/commands/commandsWithoutSignature/CD"
	"goCmd/commands/commandsWithoutSignature/Clean"
	"goCmd/debug"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strings"
)

func CMD() {
	//attempt := 0

	utils.SystemInformation()
	isWorking := true
	reader := bufio.NewReader(os.Stdin)

	prompt := ""

	for isWorking {
		if utils.IsHidden() {
			fmt.Println("You are BLOCKED!!!")
			return
		}

		dir, _ := os.Getwd()
		if prompt != "" {
			fmt.Printf("\n%s", prompt)
		} else {
			fmt.Printf("\nORPXI %s>", dir)
		}

		commandLine, _ := reader.ReadString('\n')

		commandLine = strings.TrimSpace(commandLine)
		commandParts := strings.Fields(commandLine)

		if len(commandParts) == 0 {
			continue
		}

		command := commandParts[0]
		commandArgs := commandParts[1:]
		commandLower := strings.ToLower(command)

		//isBanned := bun.UserGoCMD(command, true)
		//fmt.Println(isBanned)
		//
		//if isBanned {
		//	if attempt > 3 {
		//		bun.UserGoCMD(command, true)
		//		commandLower = "exit"
		//	} else {
		//		attempt += 1
		//	}
		//}

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

		if commandLower == "orpxihelp" {
			fmt.Println("Для получения сведений об командах наберите ORPXIHELP")
			fmt.Println("CREATE             создает новый файл")
			fmt.Println("CLEAN              очистка экрана")
			fmt.Println("CD                 смена текущего каталога")
			fmt.Println("REMOVE             удаляет файл")
			fmt.Println("READ               выводит на экран содержимое файла")
			fmt.Println("PROMPT             Изменяет ORPXI.")
			fmt.Println("PASSWORD           пароль для ORPXI.")
			fmt.Println("ORPXI              запускает ещё одну ORPXI")
			fmt.Println("SYSTEMGOCMD        вывод информации о ORPXI")
			fmt.Println("SYSTEMINFO         вывод информации о системе")
			fmt.Println("TREE               Графически отображает структуру каталогов диска или пути.")
			fmt.Println("WRITE              записывает данные в файл")
			fmt.Println("EDIT               редактирует файл")
			fmt.Println("EXIT               Выход")

			errDebug := debug.Commands(command, true)
			if errDebug != nil {
				fmt.Println(errDebug)
			}
			continue
		}

		commands := []string{"password", "promptSet", "systemgocmd", "rename", "remove", "read", "write", "create", "orpxihelp", "exit", "orpxi", "clean", "cd", "edit"}

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

		switch commandLower {
		case "password":
			Password()

		case "systemgocmd":
			utils.SystemInformation()

		case "orpxi":
			CMD()

		case "exit":
			isWorking = false

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
				continue
			}

		case "edit":
			if len(commandArgs) < 1 {
				fmt.Println("Использование: edit <файл>")
				continue
			}
			filename := commandArgs[0]
			err := Edit.File(filename)
			if err != nil {
				fmt.Println(err)
			}

		default:
			validCommand := utils.ValidCommand(commandLower, commands)
			if !validCommand {
				fmt.Printf("'%s' не является внутренней или внешней командой,\nисполняемой программой или пакетным файлом.\n", commandLine)
			}
		}
	}
}
