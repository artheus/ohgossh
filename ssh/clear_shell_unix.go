package ssh

import (
	"os"
	"os/exec"
)

func clearShell() {
	cmd := exec.Command("clear") //Windows example, its tested
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
