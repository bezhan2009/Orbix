package commands

import (
	"encoding/base64"
	"fmt"
	"goCmd/pkg/algorithms/PasswordAlgoritm"
	"goCmd/system"
	"goCmd/utils"
	"os"
	"os/exec"
	"strings"
	"time"
)

var xorMatrix = [][]int{
	{9, 10, 11, 10},
	{7, 14, 7, 14},
	{5, 18, 3, 18},
	{4, 20, 1, 20},
	{4, 41},
	{4, 41},
	{5, 39},
	{7, 35},
	{9, 31},
	{11, 27},
	{13, 23},
	{15, 19},
	{17, 15},
	{19, 11},
	{21, 7},
	{23, 3},
	{24, 1},
}

var cipherKeys = []int{
	392, 400, 440, 400, 523, 400, 392, 400,
	440, 400, 523, 400, 587, 800,
	523, 400, 587, 400, 659, 400, 523, 400,
	440, 400, 392, 400, 440, 800,
	587, 400, 659, 400, 784, 400, 587, 400,
	659, 400, 784, 400, 880, 800,
	784, 400, 880, 400, 987, 400, 784, 400,
	659, 400, 587, 400, 659, 800,
	523, 200, 659, 200, 784, 600,
}

func sysCall(freq, duration int) {
	exec.Command("powershell", "-c", fmt.Sprintf("[System.Console]::Beep(%d, %d)", freq, duration)).Run()
}

func flushBuffer() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func initSequence() {
	for i := 0; i < 3; i++ {
		sysCall(150, 150)
		time.Sleep(100 * time.Millisecond)
		sysCall(150, 150)
		time.Sleep(900 * time.Millisecond)
	}
}

func renderMatrix() {
	sym := "\u2637"
	for _, row := range xorMatrix {
		line := ""
		for i, count := range row {
			if i%2 == 0 {
				line += strings.Repeat(" ", count)
			} else {
				line += strings.Repeat(sym, count)
			}
		}
		if strings.TrimSpace(line) != "" {
			fmt.Println(system.Red(line))
			time.Sleep(150 * time.Millisecond)
		}
	}
	time.Sleep(1 * time.Second)
}

func validateChecksum() {
	sysCall(880, 200)
	time.Sleep(100 * time.Millisecond)
	sysCall(1046, 200)
	time.Sleep(150 * time.Millisecond)

	for i := 0; i < len(cipherKeys); i += 2 {
		sysCall(cipherKeys[i], cipherKeys[i+1])
		time.Sleep(50 * time.Millisecond)
	}

	sysCall(523, 300)
	time.Sleep(100 * time.Millisecond)
	sysCall(659, 300)
	time.Sleep(100 * time.Millisecond)
	sysCall(784, 600)
}

func Encrypt(commandArgs []string) {
	if len(commandArgs) < 3 {
		fmt.Println(system.Yellow("Usage: encrypt <msg> <algorithm> <key> <optional: --decrypt>"))
		fmt.Println(system.Yellow("Example: encrypt mysecretpassword xor 123"))
		fmt.Println(system.Yellow("This command has a secret message for special person :) Just send after key one more argument: secret"))
		return
	}

	encrypt := true
	if len(commandArgs) > 3 && commandArgs[3] == "--decrypt" {
		encrypt = false
	}

	if len(commandArgs) > 3 && commandArgs[3] == "secret" {
		encryptWord := PasswordAlgoritm.Usage(strings.ToLower(commandArgs[0]), true)
		hashedWord := utils.HashPasswordFromUser(encryptWord)
		encryptKey := PasswordAlgoritm.Usage(strings.ToLower(commandArgs[2]), true)
		hashedKey := utils.HashPasswordFromUser(encryptKey)

		if hashedWord != "7b6d956f499948bc7d57e284b6967fb72bba852fd9d43f2c13891fe1b8cbcf1a" &&
			hashedKey != "535fa30d7e25dd8a49f1536779734ec8286108d115da5045d77f3b4185d8f790" {
			fmt.Println(system.Yellow("Okay, there's hint: Your name is the password, but you need to find the right algorithm and key to encrypt it correctly."))
			fmt.Println(system.RedBold("README.md file contains the hint for the algorithm and key."))
		} else {
			flushBuffer()
			time.Sleep(1 * time.Second)

			initSequence()
			time.Sleep(500 * time.Millisecond)

			validateChecksum()

			renderMatrix()
			time.Sleep(1 * time.Second)

			payload := "exNFWl5fElJeRFNKQRNeXERWEkpdRh4TX0oSX1tHRl9XE1FcXVBZWlcTCBo="
			decodedMsg, err := XorDecrypt(payload, []byte(commandArgs[2]))

			fmt.Println()
			if err == nil {
				utils.AnimatedPrintLong(decodedMsg, "green")
			} else {
				fmt.Println(system.Red("System integrity check failed."))
			}
			fmt.Println("\n")

			return
		}
	}

	if strings.ToLower(commandArgs[1]) == "xor" {
		var res string
		var err error

		if encrypt {
			res, err = XorEncrypt(commandArgs[0], []byte(commandArgs[2]))
			if err != nil {
				fmt.Println(system.Red("Error: " + err.Error()))
				return
			}
			fmt.Println(system.Green("Encrypted message:\n") + system.Cyan(res))
		} else {
			res, err = XorDecrypt(commandArgs[0], []byte(commandArgs[2]))
			if err != nil {
				fmt.Println(system.Red("Error: " + err.Error()))
				return
			}
			fmt.Println(system.Red("Decrypted message:\n") + system.Red(res))
		}
	}

	if strings.ToLower(commandArgs[1]) == "klg" {
		if encrypt {
			res := PasswordAlgoritm.Klg(commandArgs[0], 10) // Assuming 10 iterations for demonstration
			fmt.Println(system.Green("Encrypted message:\n") + system.Cyan(res))
		} else {
			fmt.Println(system.Red("KLG function can't really decrypt, but you can use the same function to verify the password."))
		}
	}
}

func XorEncrypt(plaintext string, key []byte) (string, error) {
	result := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i++ {
		result[i] = plaintext[i] ^ key[i%len(key)]
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func XorDecrypt(encodedCiphertext string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("invalid base64 input: %w", err)
	}
	result := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i++ {
		result[i] = ciphertext[i] ^ key[i%len(key)]
	}
	return string(result), nil
}
