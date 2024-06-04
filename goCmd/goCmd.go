package goCmd

import (
	"bufio"
	"fmt"
	"goCmd/commands/CD"
	"goCmd/commands/Clean"
	"goCmd/commands/Create"
	"goCmd/commands/Read"
	"goCmd/commands/Remove"
	"goCmd/commands/Rename"
	"goCmd/commands/Write"
	"goCmd/debug"
	"goCmd/utils"
	"os"
	"path/filepath"
	"strings"
)

func GoCmd() {
	utils.SystemInformation()
	isWorking := true
	reader := bufio.NewReader(os.Stdin)

	prompt := ""

	for isWorking {
		dir, _ := os.Getwd()
		if prompt != "" {
			fmt.Printf("\n%s", prompt)
		} else {
			fmt.Printf("\n%s>", dir)
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

		if commandLower == "prompt" {
			fmt.Print("Enter the desired promptSet: ")
			reader := bufio.NewReader(os.Stdin)
			promptSet, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading promptSet:", err)
			} else {
				promptSet = strings.TrimSpace(promptSet)
				fmt.Printf("Prompt set to: %s\n", promptSet)
				prompt = promptSet
			}
			continue
		}

		if commandLower == "gohelp" {
			fmt.Println("Для получения сведений об командах наберите GOHELP")
			fmt.Println("CREATE             создает новый файл")
			fmt.Println("CLEAN              очистка экрана")
			fmt.Println("CD                 смена текущего каталога")
			fmt.Println("REMOVE             удаляет файл")
			fmt.Println("READ               выводит на экран содержимое файла")
			fmt.Println("PROMPT             Изменяет GoCMD.")
			fmt.Println("SYSTEMGOCMD        вывод информации о GoCMD")
			fmt.Println("SYSTEMINFO         вывод информации о системе")
			fmt.Println("GOCMD              запускает ещё одну GoCMD")
			fmt.Println("WRITE              записывает данные в файл")
			fmt.Println("EXIT               Выход")
			debug.Commands(command, true)
			continue
		}

		commands := []string{"promptSet", "systemgocmd", "rename", "remove", "read", "write", "create", "gohelp", "exit", "gocmd", "clean", "cd"}

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
		case "systemgocmd":
			utils.SystemInformation()
		case "gocmd":
			GoCmd()

		case "exit":
			isWorking = false

		case "create":
			err, name := Create.File()
			if err != nil {
				fmt.Println(err)
				debug.Commands(command, false)
			} else {
				fmt.Printf("Файл %s успешно создан!!!\n", name)
				fmt.Printf("Директория нового файла: %s\n", filepath.Join(dir, name))
				debug.Commands(command, true)
			}

		case "write":
			if len(commandArgs) < 2 {
				fmt.Println("Использование: write <файл> <данные>")
				continue
			}
			nameFileForWrite := commandArgs[0]
			data := strings.Join(commandArgs[1:], " ")

			if nameFileForWrite == "debug.txt" {
				debug.Commands(command, false)
				fmt.Println("PermissionDenied: You cannot write, delete or create a debug.txt file")
				continue
			}

			errWriting := Write.File(nameFileForWrite, data+"\n")
			if errWriting != nil {
				debug.Commands(command, false)
				fmt.Println(errWriting)
			} else {
				debug.Commands(command, true)
				fmt.Printf("Мы успешно записали данные в файл %s\n", nameFileForWrite)
			}

		case "read":
			if len(commandArgs) < 1 {
				fmt.Println("Использование: read <файл>")
				continue
			}
			nameFileForRead := commandArgs[0]

			dataRead, errReading := Read.File(nameFileForRead)
			if errReading != nil {
				debug.Commands(command, false)
				fmt.Println(errReading)
			} else {
				debug.Commands(command, true)
				_, errWrite := os.Stdout.Write(dataRead)
				if errWrite != nil {
					fmt.Println(errWrite)
				}
			}

		case "remove":
			err, name := Remove.File()
			if err != nil {
				debug.Commands(command, false)
				fmt.Println(err)
			} else {
				debug.Commands(command, true)
				fmt.Printf("Файл %s успешно удален!!!\n", name)
			}

		case "rename":
			errRename := Rename.Rename()
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

		default:
			validCommand := utils.ValidCommand(commandLower, commands)
			if !validCommand {
				fmt.Printf("'%s' не является внутренней или внешней командой,\nисполняемой программой или пакетным файлом.\n", commandLine)
			}
		}
	}
}
