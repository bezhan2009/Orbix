package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"goCmd/system"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// GitHub API URL для получения последнего коммита ветки main
const repoURL = "https://api.github.com/repos/bezhan2009/Orbix/commits/main"

// Commit Структура для обработки данных о коммите из API GitHub
type Commit struct {
	SHA string `json:"sha"`
}

func getLatestRemoteCommit() (string, error) {
	resp, err := http.Get(repoURL)
	if err != nil {
		return "", fmt.Errorf("error sending the request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server response error: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response: %v", err)
	}

	var commit Commit
	if err = json.Unmarshal(body, &commit); err != nil {
		return "", fmt.Errorf("error parsing JSON: %v", err)
	}

	return commit.SHA, nil
}

func getLocalCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error receiving a local commit: %v", err)
	}
	return string(out), nil
}

func CheckForUpdates() error {
	// Получаем последний коммит с GitHub
	remoteCommit, err := getLatestRemoteCommit()
	if err != nil {
		return err
	}

	// Получаем локальный коммит
	localCommit, err := getLocalCommit()
	if err != nil {
		return err
	}

	colors := system.GetColorsMap()

	// Сравниваем коммиты
	if strings.TrimSpace(remoteCommit) != strings.TrimSpace(localCommit) {
		var download string
		fmt.Println(colors["cyan"]("New updates are available."))
		fmt.Print(colors["cyan"]("download updates [Y/n]:"))
		reader := bufio.NewReader(os.Stdin)

		download, _ = reader.ReadString('\n')
		download = strings.TrimSpace(download)
		if strings.ToLower(download) == "y" {
			command := []string{"git", "pull", "origin", "main"}
			err = ExternalCommand(command)
			if err != nil {
				return err
			}

			commandRestart := []string{"restart.exe"}
			err = ExternalCommand(commandRestart)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println(colors["green"]("There are no updates."))
	}

	return nil
}
