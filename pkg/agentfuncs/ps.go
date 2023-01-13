package agentfuncs

// https://github.com/BishopFox/sliver/blob/master/implant/sliver/ps/ps.go

import (
	"os"
)

func Ps(){
	// TODO
}

// Kill finds a process from a PID and terminates it.
func Kill(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := p.Kill(); err != nil {
		return err
	}
	_, err = p.Wait()
	return err
}