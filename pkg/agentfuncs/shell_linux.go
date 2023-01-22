//go:build linux
// +build linux

package agentfuncs

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

// Execute via system shell so we have access to shell internal commands
func Shell(cmd []string) string {
    res, err := exec.Command("/bin/sh", append([]string{"-c"}, strings.Join(cmd, " "))...).CombinedOutput()
    if (err != nil) {
        return string(string(res) + "\n[!] Golang Error: " + err.Error()) 
    }
    return string(res)
}

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

    var sb strings.Builder
    var line string
    var stat *syscall.Stat_t
	for _, file := range fileInfo {
		stat = file.Sys().(*syscall.Stat_t)
		uid := int(stat.Uid)
		gid := int(stat.Gid)
		fUser, _ := user.LookupId(strconv.Itoa(uid))
		fGroup, _ := user.LookupGroupId(strconv.Itoa(gid)) 
		// TODO: Account for symlinks
		if file.IsDir() {
			// drw-rw-rw 1000 1000 date name
			line = fmt.Sprintf("d%s\t%s\t%s\t%s", file.Mode().Perm().String()[1:], fUser, fGroup, file.Name())
		} else {
			line = fmt.Sprintf("%s\t%s\t%s\t%s", file.Mode().Perm().String(), fUser, fGroup, file.Name())
		}
		sb.WriteString(line + "\n")
	}

    return sb.String()
}