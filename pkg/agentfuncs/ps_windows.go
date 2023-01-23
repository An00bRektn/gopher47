//go:build windows
// +build windows
package agentfuncs

import (
	"fmt"
	"os"
	"os/user"
	"text/tabwriter"
	"strings"
	"github.com/elastic/go-sysinfo"
)

// TODO: This works very well on windows, not so well on Linux
// 		 maybe we use a different package?
// 		 specifically, it only works on linux with sudo perms
func Ps() string {
	procs, err := sysinfo.Processes()
	var sb strings.Builder
	if err != nil {
		return "[!] Error: " + err.Error()
	}
	writer := tabwriter.NewWriter(&sb, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "USER\tPID\tPPID\tSTART\tCOMMAND")
	fmt.Fprintln(writer, "=====\t====\t=====\t=====\t========")

	var info string
	for _, proc := range procs {
		procInfo, err := proc.Info()
		if err != nil {
			return "[!] Error: " + err.Error()
		}

		userInfo, err := proc.User()
		if err != nil {
			return "[!] Error: " + err.Error()
		}

		uid := userInfo.UID
		userobj, err := user.LookupId(uid)
		var username string
		if err != nil {
			username = uid
		} else {
			username = userobj.Username
		}

		info = fmt.Sprintf("%s\t%d\t%d\t%s\t%s", username, procInfo.PID, procInfo.PPID, procInfo.StartTime.Format("2:00"), procInfo.Exe)
		fmt.Fprintln(writer, info)
	}
	writer.Flush()
	return sb.String()
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