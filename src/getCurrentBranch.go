package src

import (
	"bytes"
	"os/exec"
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
