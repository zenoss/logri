package logri

type LogriConfig struct {
	Loggers map[string]LoggerConfig `json:"loggers"`
}

type LoggerConfig struct {
	Level string `json:"level"`
}
