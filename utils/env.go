package utils

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}
	token := os.Getenv(key)
	if token == "" {
		return "", errors.New("missing api key")
	}

	return token, nil
}
