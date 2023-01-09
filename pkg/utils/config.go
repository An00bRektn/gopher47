package utils

import (
	"bytes"
	"encoding/json"
	"strings"
)

type Config struct {
	Url string
	SleepTime int
	JitterRange int
}

func GetConfig() Config {
	config := Config{
		Url: "http://127.0.0.1:80",//"{{URL}}",
		SleepTime: 5,//{{SleepTime}},
		JitterRange: 100,//{{JitterRange}}
	}
	return config
}

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