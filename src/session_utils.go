package src

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// CheckGit checks if Git is installed on the system
func CheckGit() bool {
	// Попробуем выполнить команду "git --version"
	_, err := exec.LookPath("git")
	if err != nil {
		log.Println("Git is not installed.")
		return false
	}
	return true
}

func GetCurrentGitBranch() (string, error) {
	if !GitCheck {
		ErrGitNotInstalled := errors.New("ErrGitNotInstalled")
		return "", ErrGitNotInstalled
	}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	branch := out.String()
	return branch[:len(branch)-1], nil
}

func Getwd() (wd string) {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		fmt.Println(yellow("WARNING: Some commands may not work because the Getwd function failed with an error"))
		fmt.Println(red(err))
	}

	wd = strings.TrimSpace(wd)

	return
}
