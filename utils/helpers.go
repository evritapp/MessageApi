package utils

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
)

func CheckModel(sms int) (bool, error) {
	if sms > 1 {
		return false, errors.New("SendingType not correct")
	}
	return true, nil
}

func InitEnvVars() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file", err)
		return err
	}
	return nil
}
