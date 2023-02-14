package config

import (
	"github.com/spf13/pflag"
	"os"
	"strconv"
	"time"
)

const (
	defaultAddress         = "localhost:8080"
	defaultAccrualAddress  = "localhost:8081"
	defaultDBConfig        = ""
	defaultSessionLifetime = time.Minute * 30
)

type AppConfig struct {
	ServerAddress   string
	AccrualAddress  string
	DbConnection    string
	SessionLifetime time.Duration
}

func InitAppConfig() *AppConfig {
	address := getAddress()
	var addressF string
	pflag.StringVarP(&addressF, "a", "a", defaultAddress, "Address of the server")

	accrualAddress := getAccrualAddress()
	var accrualAddressF string
	pflag.StringVarP(&accrualAddressF, "r", "r", defaultAccrualAddress, "Address of an accrual server")

	dbConfig := getDBConfig()
	var dbDsnF string
	pflag.StringVarP(&dbDsnF, "d", "d", defaultDBConfig, "Postgres DB URI")

	pflag.Parse()

	if address == "" {
		address = addressF
	}
	if accrualAddress == "" {
		accrualAddress = accrualAddressF
	}
	if dbConfig == "" {
		dbConfig = dbDsnF
	}

	return &AppConfig{
		ServerAddress:   address,
		AccrualAddress:  accrualAddress,
		DbConnection:    dbConfig,
		SessionLifetime: getSessionLifetime(),
	}
}

func getAddress() string {
	return os.Getenv("RUN_ADDRESS")
}

func getAccrualAddress() string {
	return os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
}
func getDBConfig() string {
	return os.Getenv("DATABASE_URI")
}
func getSessionLifetime() time.Duration {
	reportInterval := os.Getenv("SESSION_LIFETIME")
	if reportInterval == "" {
		return defaultSessionLifetime
	}
	reportIntervalVal, err := strconv.Atoi(reportInterval)
	if err != nil {
		return defaultSessionLifetime
	}
	duration := time.Minute * time.Duration(reportIntervalVal)
	return duration
}
