package config

import (
	"encoding/json"
	"fmt"
	zlog "github.com/rs/zerolog/log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName    string `json:"app_name"`
	ServerPort string `json:"server_port"`

	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBUser     string `json:"-"`
	DBPassword string `json:"-"`
	DBName     string `json:"db_name"`
	DBSSLMode  string `json:"db_sslmode"`

	DatabaseURL string `json:"-"`
}

func LoadConfig(configPath string) *Config {
	_ = godotenv.Load()

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		zlog.Fatal().Msgf("Error reading config file '%s': %v\n", configPath, err)
	}

	var cfg Config
	err = json.Unmarshal(configFile, &cfg)
	if err != nil {
		zlog.Fatal().Msgf("Error unmarshalling config JSON from '%s': %v\n", configPath, err)
	}

	dbUsername := os.Getenv("DB_USERNAME")
	if dbUsername == "" {
		zlog.Fatal().Msgf("Error: DB_USERNAME environment variable is not set.")
	}
	cfg.DBUser = dbUsername

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		zlog.Fatal().Msgf("Error: DB_PASSWORD environment variable is not set.")
	}
	cfg.DBPassword = dbPassword

	cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		strconv.Itoa(cfg.DBPort),
		cfg.DBName,
		cfg.DBSSLMode,
	)

	zlog.Info().Msgf("Configuration loaded successfully for app: %s", cfg.AppName)
	return &cfg
}
