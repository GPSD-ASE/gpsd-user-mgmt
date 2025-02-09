package config

import (
	"os"
)

type Config struct {
	PORT    string
	DB_NAME string
	DB_HOST string
	DB_PORT string
	DB_USER string
	DB_PASS string
}

func Load() *Config {
	return &Config{
		PORT:    os.Getenv("USER_PORT"),
		DB_NAME: os.Getenv("DB_NAME"),
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PASS: os.Getenv("DB_PASS"),
	}
}
