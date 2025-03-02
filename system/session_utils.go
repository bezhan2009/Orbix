package system

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var GitExists = CheckPackageExists("git")

func GetCurrentGitBranch() (string, error) {
	if !GitExists {
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
		fmt.Println(Yellow("WARNING: Some commands may not work because the Getwd function failed with an error"))
		fmt.Println(Red(err))
	}

	wd = strings.TrimSpace(wd)

	return
}

func SetGitBranch(sd *Session) {
	var errGitBranch error
	sd.GitBranch, errGitBranch = GetCurrentGitBranch()
	if errGitBranch != nil {
		sd.GitBranch = ""
	}
}

func SetPath(sd *Session) {
	sd.Path = Getwd()
}
