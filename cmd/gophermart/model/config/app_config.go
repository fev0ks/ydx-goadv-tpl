package config

import (
	"github.com/spf13/pflag"
	"os"
)

const (
	defaultAddress = "localhost:8080"
)

type AppConfig struct {
	ServerAddress string
}

func InitAppConfig() *AppConfig {
	address := getAddress()
	var addressF string
	pflag.StringVarP(&addressF, "a", "a", defaultAddress, "Address of the server")

	pflag.Parse()
	if address == "" {
		address = addressF
	}
	return &AppConfig{
		ServerAddress: address,
	}
}

func getAddress() string {
	return os.Getenv("ADDRESS")
}
