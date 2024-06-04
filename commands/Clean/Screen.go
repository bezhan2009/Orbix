package Clean

import (
	"os"
	"os/exec"
)

func Screen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
