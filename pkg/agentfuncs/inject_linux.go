// +build linux

package agentfuncs

// TODO: Consider possibly doing that trick to run shellcode from the commandline?
//		 Maybe just do what OffensiveNotion does and make it a dropper?
func SelfInject(shellcodeHex string) string {
	return "[!] Not currently supported for Linux!"
}