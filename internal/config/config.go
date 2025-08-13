package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server Server
}

type Server struct {
	Host        string        `env:"host" env-default:"localhost"`
	Port        string        `env:"port" env-default:"8080"`
	Timeout     time.Duration `env:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `env:"idle_timeout" env-default:"30s"`
}

const configPath = "config/local.env"

func MustLoad() *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not exist: " + configPath)
	}

	if err := godotenv.Load(configPath); err != nil {
		log.Fatalf("cannot load env file: %s", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config: " + err.Error())
	}

	return &cfg
}
