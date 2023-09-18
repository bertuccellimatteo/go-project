package env

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func ValidateEnv() error {

	err := godotenv.Load(".env")

	if err != nil {
		return err
	}

	{
		_, isFound := os.LookupEnv("DATABASE_URL")
		if !isFound {
			return errors.New("DATABASE_URL env variable not found")
		}
	}

	{
		var _, isFound = os.LookupEnv("ACCESS_TOKEN_SECRET")
		if !isFound {
			return errors.New("ACCESS_TOKEN_SECRET env variable not found")
		}
	}

	return nil
}
