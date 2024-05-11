package config

type LoggerConfig struct {
	Level    string `json:"level" envconfig:"LOGGER_LEVEL"`
	Encoding string `json:"encoding" envconfig:"LOGGER_ENCODING"`
}