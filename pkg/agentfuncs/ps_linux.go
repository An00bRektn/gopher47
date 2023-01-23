//go:build linux
// +build linux

package agentfuncs

// Reference: https://github.com/BishopFox/sliver/blob/master/implant/sliver/ps/ps_linux.go
// 			  https://github.com/mitchellh/go-ps
import (
	"fmt"
	"io"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
)

type UnixProcess struct {
	UID		int
	User	string
	PID		int
	PPID	int
	Command string	
}

func getCmdLine(pid int) (string, error) {
	path := fmt.Sprintf("/proc/%d/cmdline", pid)
	dat, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	args := strings.Split(string(dat), "\x00")
	return strings.Join(args, " "), nil
}

func getProcUID(pid int) (int, error) {
	path := fmt.Sprintf("/proc/%d/task", pid)
	fd, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	fileStat := &syscall.Stat_t{}
	err = syscall.Fstat(int(fd.Fd()), fileStat)
	if err != nil {
		return 0, err
	}
	return int(fileStat.Uid), nil
}

func getProcUser(pid int) (string, error) {
	uid, err := getProcUID(pid)
	if err != nil {
		return "", err
	}
	usr, err := user.LookupId(fmt.Sprintf("%d", uid))
	if err != nil {
		return fmt.Sprintf("%d", uid), err
	}
	return usr.Username, nil
}

func getProcPPID(pid int) (int, error) {
	path := fmt.Sprintf("/proc/%d/stat", pid)
	dat, err := os.ReadFile(path)
	if err != nil {
		return -1, err
	}
	ppid, err := strconv.Atoi(strings.Split(string(dat), " ")[3])
	if err != nil {
		return -1, err
	}
	return ppid, nil
}

func newUnixProcess(pid int) (*UnixProcess, error){
	ppid, err := getProcPPID(pid)
	if err != nil {
		ppid = -1
	}
	uid, err := getProcUID(pid)
	if err != nil {
		uid = -1
	}
	username, err := getProcUser(pid)
	if err != nil {
		if uid < 0 {
			username = ""
		} else {
			username = fmt.Sprintf("%d", uid)
		}
	}
	cmd, err := getCmdLine(pid)
	if err != nil {
		cmd = ""
	}

	proc := &UnixProcess{
		UID:		uid,
		User:		username,
		PID:		pid,
		PPID:		ppid,
		Command:	cmd,
	}
	return proc, err
}

func processes() ([]UnixProcess, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	procs := []UnixProcess{}
	for {
		names, err := d.Readdirnames(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, name := range names {
			if name[0] < '0' || name[0] > '9' {
				continue
			}
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}

			p, err := newUnixProcess(int(pid))
			if err != nil {
				continue
			}

			procs = append(procs, *p)
		}
	}

	return procs, nil
}

func Ps() string {
	procs, err := processes()
	if err != nil {
		return "[!] Error: " + err.Error()
	}
	var sb strings.Builder
	writer := tabwriter.NewWriter(&sb, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, "USER\tPID\tPPID\tCOMMAND")
	fmt.Fprintln(writer, "=====\t====\t=====\t========")

	var info string
	for _, proc := range procs {
		info = fmt.Sprintf("%s\t%d\t%d\t%s", proc.User, proc.PID, proc.PPID, proc.Command)
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