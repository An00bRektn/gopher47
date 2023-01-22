//go:build windows
// +build windows

package agentfuncs

import (
	"os/exec"
    "os"
    "strings"
	"syscall"
)

func Shell(cmd []string) string {
    c := exec.Command("cmd.exe", append([]string{"/c"}, strings.Join(cmd, " "))...)
	c.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	res, err := c.CombinedOutput()
    if (err != nil) {
        return string(res) + "\n[!] Golang Error: " + err.Error()
    }
    return string(res)
}

func PowerShell(cmd []string) string {
    c := exec.Command("powershell.exe", append([]string{"-Command"}, strings.Join(cmd, " "))...)
	c.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	res, err := c.CombinedOutput()
    if (err != nil) {
        return string(res) + "\n[!] Golang Error: " + err.Error()
    }
    return string(res)
}

// TODO: Make a Windows version that looks nice
func Ls(dir string) string {
    // Read Directory Listing
	f, err := os.Open(dir)
	if err != nil {
		return "[!] Error: " + err.Error()
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()

	if err != nil {
		return "[!] Error: " + err.Error()
	}

    // TODO: Print a nice output with info and stuff
    var sb strings.Builder
	for _, file := range fileInfo {
		sb.WriteString(file.Name() + "\n")
	}

    return sb.String()
}