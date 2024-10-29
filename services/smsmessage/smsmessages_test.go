package smsmessage_test

import (
	"errors"
	"testing"

	"messageapi.e-vrit.co.il/services/smsmessage"
)

func TestSmsModel(t *testing.T) {

	isms, err := 1, errors.New("WRONG MESSAGE")

	if err != nil && isms==1 {
		t.Error("test error", err)
	}

}
func TestSendSms(t *testing.T) {


	_, err := smsmessage.CheckModel(1)

	if err != nil {
		t.Error("test error", err)
	}

}