package smsmessage_test

import (
	"errors"
	"fmt"
	"testing"

	"messageapi.e-vrit.co.il/services/smsmessage"
	"messageapi.e-vrit.co.il/services/smsmessage/flashy"
	"messageapi.e-vrit.co.il/services/smsmessage/models"
)

func TestSmsModel(t *testing.T) {

	isms, err := 1, errors.New("WRONG MESSAGE")

	if err != nil && isms == 1 {
		t.Error("test error", err)
	}

}
func TestSendSms(t *testing.T) {

	var sms models.SmsModel // a == Student{"", 0}
	sms.Message = "12345"
	sms.ReciverPhoneNumber = "0587198145"
	sms.SenderName = "e-vrit"
	sms.SendingType = 1

	var err error
	var isms smsmessage.ISms

	isms, err = flashy.NewFlashySmsModel()
	fmt.Println("isms: ", isms, "err: ", err)
	if err != nil {
		t.Error("test error", err)

	}
	// isms.SendSms(sms)
	fmt.Println(isms)

	// _, err := smsmessage.CheckModel(2)

	// if err != nil {
	// 	t.Error("test error", err)
	// }

}
