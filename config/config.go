package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env       string `json:"env" envconfig:"APP_ENV"`
	BotToken  string `json:"BOT_TOKEN" envconfig:"BOT_TOKEN"`
	OpenAIKey string `json:"OPENAI_KEY" envconfig:"OPENAI_KEY"`
	DB *DBConfig
	Logger *LoggerConfig
}

func InitConfig() (*Config, error) {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
