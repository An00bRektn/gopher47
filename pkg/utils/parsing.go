package utils

import (
	"bytes"
	"encoding/json"
	"strings"
)

/* Removes null bytes at end of string, then whitespace, and casts to string */
func Strip(dirtyString string) string {
	return strings.TrimSpace(string(bytes.Trim([]byte(dirtyString), "\x00")))
}

// https://stackoverflow.com/questions/51691901/how-do-you-escape-characters-within-a-string-json-format
/*
	Escapes characters that would otherwise break JSON parsing
	Input: i string - input string to be escaped
	Output: escaped string
*/
func JsonEscape(i string) string {
    b, err := json.Marshal(i)
    if err != nil {
        panic(err)
    }
    // Trim the beginning and trailing " character
    return string(b[1:len(b)-1])
}

func FindNotLoopback(ips []string) string {
	var cleaned string
	for _, ip := range ips {
		cleaned = strings.Split(ip, "/")[0]
		// apologies if you're an IPv6 user
		// but this is what you change if you need that
		if cleaned != "127.0.0.1" && !strings.Contains(cleaned, ":") {
			return cleaned
		}
	}
	return "0.0.0.0"
}