package src

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetCurrentGitBranch() (string, error) {
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
