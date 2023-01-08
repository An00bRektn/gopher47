package functions

import (
	"os/exec"
)

func Shell(cmd string) string {
	res, _ := exec.Command(cmd).Output()
	return string(res) + "\n"
}
