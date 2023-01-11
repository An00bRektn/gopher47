package agentfuncs

import (
	"os/exec"
    "os"
    "strings"
)

func Shell(cmd []string) string {
    res, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
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

    // TODO: Print a nice output with info and stuff
    var sb strings.Builder
	for _, file := range fileInfo {
		sb.WriteString(file.Name() + "\n")
	}

    return sb.String()
}