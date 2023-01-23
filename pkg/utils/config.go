package utils
// echo 'goodeveninggopher47' | sha256sum

import (
	"reflect"
)

type Config struct {
	Url string
	IsSecure bool
	UserAgent string
	SleepTime int
	JitterRange int
	TimeoutThreshold int
}

/* 
	Returns config to agent. Make modifications here.
	We're defining the config here so it's more "malleable",
	although this isn't a one-to-one
*/
func GetConfig() Config {
	config := Config{
		Url: "http://127.0.0.1:80/",
		IsSecure: false,
		UserAgent: "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
		SleepTime: 10,
		JitterRange: 100,
		TimeoutThreshold: 4,
	}
	return config
}

// Message - Fake message for embedding canaries
type Message struct {
	Command string `c2:"cb701b6f0a2f55e3c269f5dde3f4ba25f55be6e65add1657b6843430bf1a4940"`
}

// never obfuscate the Message type
var _ = reflect.TypeOf(Message{})