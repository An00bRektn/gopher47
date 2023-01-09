package utils

type Config struct {
	Url string
	SleepTime int
	JitterRange int
}

/* Returns config to agent. Make modifications here. */
func GetConfig() Config {
	config := Config{
		Url: "http://127.0.0.1:80",
		SleepTime: 5,
		JitterRange: 100,
	}
	return config
}