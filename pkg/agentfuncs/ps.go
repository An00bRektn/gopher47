package functions

// https://github.com/BishopFox/sliver/blob/master/implant/sliver/ps/ps.go

import (
	"fmt"
	"os"
)

func Ps(){

}

// Kill finds a process from a PID and terminates it.
func Kill(pid int) error {
	fmt.Println(pid)
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()	
}