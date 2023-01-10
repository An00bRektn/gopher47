package functions

import (
	"os/exec"
)

func Shell(cmd []string) string {
    res, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
    if (err != nil) {
        return string(string(res) + "\n[!] Golang Error: " + err.Error()) 
    }
    return string(res)
}
