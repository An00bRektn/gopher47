package functions

import (
	//"log"
	"os/exec"
)

func Shell(cmd []string) string {
    res, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
    if (err != nil) {
		//log.Fatal(err)
        return string(string(res) + "\n[!] Golang Error:" + err.Error()) 
    }
    return string(res)
}
