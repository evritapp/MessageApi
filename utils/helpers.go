package utils

import (
	"errors"
)



func CheckModel(sms int) (bool, error) {
	if sms > 1 {
		return false, errors.New("SendingType not correct")
	}
	return true, nil
}
