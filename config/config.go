package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Configuration struct {
	DBConn     string
	ServerAddr string
	JwtSign    string
}

var Env string
var Config Configuration

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		if value == "" {
			return fallback
		}
		return value
	}
	return fallback
}

func LoadEnv() (config Configuration, err error) {
	env := GetEnv("ENV", "local")
	if err := godotenv.Load(".env"); err != nil {
		godotenv.Load("../.env")
	}

	allowedDomainStr := os.Getenv("ALLOWED_DOMAINS")
	allowedDomains := []string{}
	if len(allowedDomainStr) > 0 {
		allowedDomains = strings.Split(allowedDomainStr, ",")
	}

	for i := range allowedDomains {
		allowedDomains[i] = strings.TrimPrefix(allowedDomains[i], " ")
	}

	config = Configuration{
		DBConn:     os.Getenv("DB_CONN"),
		ServerAddr: os.Getenv("SERVER_ADDR"),
		JwtSign:    os.Getenv("JWT_SIGN"),
	}
	Env = env
	Config = config

	return config, nil
}
