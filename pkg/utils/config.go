package utils

type Config struct {
	Url string
	SleepTime int
	JitterRange int
	TimeoutThreshold int
}

/* Returns config to agent. Make modifications here. */
func GetConfig() Config {
	config := Config{
		Url: "http://10.10.69.24:8080/",
		SleepTime: 15,
		JitterRange: 100,
		TimeoutThreshold: 4,
	}
	return config
}