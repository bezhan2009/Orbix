package src

import (
	"bufio"
	"errors"
	"fmt"
	_chan "goCmd/chan"
	"goCmd/cmd/dirInfo"
	"goCmd/src/service"
	"goCmd/src/user"
	"goCmd/structs"
	"goCmd/system"
	"os"
	"strings"
)

func DefineUser(commandInput string,
	rebooted structs.RebootedData,
	sessionData *system.AppState) (string, error) {
	var username string

	// Check if password directory is empty once and handle errors here
	isEmpty, err := user.IsPasswordDirectoryEmpty()
	if err != nil {
		service.AnimatedPrint(fmt.Sprintf("Error checking password directory: %s\n", err.Error()), "red")
		return "", errors.New("ErrCheckPasswordDirectory")
	}

	if strings.TrimSpace(rebooted.Username) != "" {
		username = strings.TrimSpace(rebooted.Username)
	} else if !isEmpty && commandInput == "" {
		dir, _ := os.Getwd()
		OrbixUser := dirInfo.CmdUser(&dir)

		nameUser, isSuccess, errUser := user.CheckUser(OrbixUser, sessionData)
		if !isSuccess {
			if _chan.UserStatusAuth {
				system.Unauthorized = false
				_chan.UpdateChan("system__user")
			} else {
				system.Unauthorized = true
			}

			return "", errUser
		}

		system.Unauthorized = false
		username = nameUser
		if username != OrbixUser {
			initializeRunningFile(username)
		}

		if OrbixUser == username {
			sessionData.IsAdmin = true
			sessionData.User = OrbixUser
		} else {
			sessionData.IsAdmin = false
			sessionData.User = username
		}
	}

	return username, nil
}

func OrbixUser(commandInput string,
	echo bool,
	rebooted *structs.RebootedData,
	SD *system.AppState,
	ExecLtCommand func(commandInput string)) (LoopData structs.OrbixLoopData, LoadUserConfigsFn func(echo bool) error) {
	LoadUserConfigsFn = LoadConfigs

	if UsingForLT(commandInput) {

		// Load User Configs
		_ = LoadConfigs(true)

		ExecLtCommand(commandInput)

		isWorking := false
		isPermission := false
		RestartAfterInit := false

		return structs.OrbixLoopData{
			IsWorking:        &isWorking,
			IsPermission:     &isPermission,
			Username:         "",
			SessionData:      &system.AppState{},
			RestartAfterInit: &RestartAfterInit,
		}, LoadUserConfigsFn
	}

	RestartAfterInit := false

	sessionData := InitOrbixFn(&RestartAfterInit,
		echo,
		commandInput,
		*rebooted,
		SD)

	isWorking := true
	isPermission := true
	if commandInput != "" {
		isPermission = false
	}

	username, err := DefineUser(commandInput,
		*rebooted,
		sessionData)
	if err != nil {
		isWorking = false
		isPermission = false
		RestartAfterInit = false

		return structs.OrbixLoopData{
			IsWorking:        &isWorking,
			IsPermission:     &isPermission,
			Username:         "",
			SessionData:      &system.AppState{},
			RestartAfterInit: &RestartAfterInit,
		}, LoadUserConfigsFn
	}

	// Load User Configs
	_ = LoadConfigs(true)

	if username != "" {
		system.EditableVars["user"] = &username
	}

	return structs.OrbixLoopData{
		IsWorking:        &isWorking,
		IsPermission:     &isPermission,
		Username:         username,
		SessionData:      sessionData,
		RestartAfterInit: &RestartAfterInit,
	}, nil
}

func GetUserNickname() string {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(system.Magenta("\nYour name is Empty!!!\nEnter nickname: "))
		nickname, _ := reader.ReadString('\n')

		nickname = strings.TrimSpace(nickname)

		if nickname == "" {
			continue
		}

		return nickname
	}
}
