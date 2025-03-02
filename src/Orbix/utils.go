package Orbix

import (
	"fmt"
	"goCmd/cmd/commands"
	"goCmd/src"
	"goCmd/src/user"
	"goCmd/structs"
	"goCmd/system"
	"log"
	"os"
	"strings"
)

func RecoverFromThePanic(commandInput string,
	r any,
	echo bool,
	SD *system.AppState,
	OrbixLoopData structs.OrbixLoopData,
	LoadUserConfigsFn func(echo bool) error) {
	PanicText := fmt.Sprintf("Panic recovered: %v", r)
	fmt.Printf("\n%s\n", system.Red(PanicText))

	if system.RebootAttempts > system.MaxRetryAttempts {
		fmt.Println(system.Red("Max retry attempts reached. Exiting..."))
		log.Println("Max retry attempts reached. Exiting...")
		os.Exit(1)
	}

	system.RebootAttempts += 1

	fmt.Println(system.Yellow("Recovering from panic"))

	log.Printf("Panic recovered: %v", r)

	var reboot = structs.RebootedData{
		Username:          system.UserName,
		Recover:           r,
		Prefix:            system.Prefix,
		LoopData:          OrbixLoopData,
		LoadUserConfigsFn: LoadUserConfigsFn,
	}

	Orbix(strings.TrimSpace(commandInput),
		echo,
		reboot,
		SD)
}

func RestartAfterInitFn(sessionData *system.AppState,
	rebooted structs.RebootedData,
	prefix,
	username string,
	echo bool) {
	sessionData.User = username
	rebooted.Prefix = prefix
	if len(os.Args) > 1 {
		return
	}

	Orbix("",
		echo,
		rebooted,
		sessionData)
}

func handlePanic(commandInput string,
	echo bool,
	SD *system.AppState,
	OrbixLoopData structs.OrbixLoopData,
	LoadUserConfigsFn func(echo bool) error,
	r any) {
	system.Unauthorized = true
	if r != nil {
		user.DeleteUserFromRunningFile(system.User)

		system.OrbixRecovering = true
		RecoverFromThePanic(
			commandInput,
			r,
			echo,
			SD,
			OrbixLoopData,
			LoadUserConfigsFn)
	}
}

func setupOutputRedirect(echo bool) (originalStdout, originalStderr *os.File) {
	originalStdout, originalStderr = os.Stdout, os.Stderr
	devNull, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0666)
	defer func() {
		err := devNull.Close()
		if err != nil {
			return
		}
	}()

	if echo {
		os.Stdout, os.Stderr = originalStdout, originalStderr
	} else {
		os.Stdout, os.Stderr = devNull, devNull
	}

	return
}

func EndOfSessions(originalStdout, originalStderr *os.File,
	session *system.Session,
	sessionData *system.AppState,
	prefix string) {
	system.CntLaunchedOrbixes--

	// Restore original outputs
	os.Stdout, os.Stderr = originalStdout, originalStderr

	system.PreviousSessionPath = session.Path
	session, _ = sessionData.GetSession(system.PreviousSessionPrefix)

	if strings.TrimSpace(session.Path) != "" {
		if err := commands.ChangeDirectory(session.Path); err != nil {
			fmt.Println(system.Red("Error changing directory:", err))
		}
	}

	sessionData.DeleteSession(prefix)

	system.OrbixWorking = false
	src.UnknownCommandsCounter = 0
}

func RecoverAndRestore(rebooted *structs.RebootedData) {
	if rebooted.Recover != nil {
		RecoverText := fmt.Sprintf("Successfully recovered from the panic: %v", rebooted.Recover)
		fmt.Printf("\n%s\n", system.Green(RecoverText))
		rebooted.Recover = nil
	}

	system.OrbixRecovering = false
}
