package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	ServerAddress   string
	DatabaseAddress string
	AccrualAddress  string
}

func parseFlags() (string, string, string) {
	serverAddressFlag := flag.String("a", "localhost:8080", "адрес и порт запуска сервера")
	databaseAddressFlag := flag.String("d", "", "адрес подключения к базе данных")
	accrualAddressFlag := flag.String("r", "", "адрес системы расчёта начислений")

	flag.Parse()
	return *serverAddressFlag, *databaseAddressFlag, *accrualAddressFlag
}

func parseEnv() (string, string, string) {
	envServerAddress := os.Getenv("RUN_ADDRESS")
	envDatabaseAddress := os.Getenv("DATABASE_URI")
	envAccrualAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	return envServerAddress, envDatabaseAddress, envAccrualAddress
}

// NewConfig создает новый экземпляр конфигурации приложения на основе флагов командной строки и переменных окружения
func NewConfig() *Config {
	serverAddressFlag, databaseAddressFlag, accrualAddressFlag := parseFlags()
	envServerAddress, envDatabaseAddress, envAccrualAddress := parseEnv()

	cfg := &Config{}

	if envServerAddress != "" {
		cfg.ServerAddress = envServerAddress
	} else {
		cfg.ServerAddress = serverAddressFlag
	}

	if envDatabaseAddress != "" {
		cfg.DatabaseAddress = envDatabaseAddress
	} else {
		cfg.DatabaseAddress = databaseAddressFlag
	}

	if envAccrualAddress != "" {
		cfg.AccrualAddress = envAccrualAddress
	} else {
		cfg.AccrualAddress = accrualAddressFlag
	}

	cfg.ServerAddress = strings.TrimPrefix(cfg.ServerAddress, "http://")
	parts := strings.Split(cfg.ServerAddress, ":")
	if len(parts) > 1 && parts[0] == "" {
		cfg.ServerAddress = "localhost:" + parts[1]
	}

	return cfg
}
