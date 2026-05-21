package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	logger *log.Logger
	cfg    *Config
)

func init() {
	loadDotEnv()

	var confName string
	env := os.Getenv("ENV")

	if env == "container" {
		confName = "container"
	} else {
		confName = "local"
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "."
	}

	// init log system
	logger = log.Default()
	logger.SetOutput(os.Stdout)

	viper.SetConfigName(confName) // Name of the config file (without extension)
	viper.SetConfigType("yaml")   // Config file type
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("config")

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Printf("Warning: Config file not found, using default values: %v", err)
		} else {
			logger.Fatalf("Error reading config file: %v", err)
		}
	}

	// Initialize variables by unmarshaling into the struct
	cfg = new(Config)
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	if err := applyDBEnv(cfg); err != nil {
		if isTestBinary() {
			logger.Printf("Warning: database configuration skipped for tests: %v", err)
		} else {
			logger.Fatalf("database configuration error: %v", err)
		}
	}
}

func loadDotEnv() {
	// Load uses existing OS environment variables as the source of truth and only
	// fills missing values from .env when the file exists.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: unable to load .env file: %v", err)
	}
}

func isTestBinary() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

func applyDBEnv(cfg *Config) error {
	cfg.DBConfig = DBConfig{
		Host:     strings.TrimSpace(os.Getenv("POSTGRES_HOST")),
		Port:     strings.TrimSpace(os.Getenv("POSTGRES_PORT")),
		User:     strings.TrimSpace(os.Getenv("POSTGRES_USER")),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     strings.TrimSpace(os.Getenv("POSTGRES_DB")),
		SSLMode:  strings.TrimSpace(os.Getenv("POSTGRES_SSLMODE")),
	}

	if cfg.DBConfig.SSLMode == "" {
		cfg.DBConfig.SSLMode = "disable"
	}

	missing := make([]string, 0)
	for key, value := range map[string]string{
		"POSTGRES_HOST":     cfg.DBConfig.Host,
		"POSTGRES_PORT":     cfg.DBConfig.Port,
		"POSTGRES_USER":     cfg.DBConfig.User,
		"POSTGRES_PASSWORD": cfg.DBConfig.Password,
		"POSTGRES_DB":       cfg.DBConfig.Name,
	} {
		if value == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s. Create .env or export them in the OS environment", strings.Join(missing, ", "))
	}

	return nil
}

func GetLogger() *log.Logger {
	return logger
}

func GetConfig() *Config {
	return cfg
}
