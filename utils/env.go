package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) (string, error) {

	if err := godotenv.Load(); err != nil {
		return "", err
	}
	token := os.Getenv(key)

	return token, nil
}
