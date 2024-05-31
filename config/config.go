package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Smtp struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
}

func LoadConfig() *Smtp {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("smtpConfig path is empty")
	}
	var smtpConfig Smtp
	err := cleanenv.ReadConfig(path, &smtpConfig)
	if err != nil {
		panic("failed to read smtpConfig file " + path)
	}

	return &smtpConfig
}
