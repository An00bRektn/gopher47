package utils

type Config struct {
	Url string
	IsSecure bool
	UserAgent string
	SleepTime int
	JitterRange int
	TimeoutThreshold int
}

/* Returns config to agent. Make modifications here. */
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