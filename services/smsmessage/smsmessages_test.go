package smsmessage_test

import (
	"fmt"
	"os"
	"testing"

	"messageapi.e-vrit.co.il/services/smsmessage"
	"messageapi.e-vrit.co.il/services/smsmessage/flashy"
	"messageapi.e-vrit.co.il/services/smsmessage/models"
)

// func loadEnv() error {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		return fmt.Errorf("error loading .env file: %v", err)
// 	}

// 	env := os.Getenv("GO_ENV")
// 	if env == "" {
// 		env = "development" // Default to development if not set
// 	}
// 	return godotenv.Overload(fmt.Sprintf(".env.%s", env))
// }

func funcUsingEnvVar() string {
	return os.Getenv("GO_ENV")
}

func setEnvVarForTesting(t *testing.T) {
	t.Setenv("GO_ENV", "dev")
}

func TestSettingEnvVar(t *testing.T) {
	setEnvVarForTesting(t)
	if funcUsingEnvVar() != "dev" {
		t.Fatal("GO_ENV not set to dev")
	}
}

// func TestSmsModel(t *testing.T) {

// 	isms, err := 1, errors.New("WRONG MESSAGE")

// 	if err != nil && isms == 1 {
// 		t.Error("test error", err)
// 	}

// }
func TestSendSms(t *testing.T) {

	// if err := TestSettingEnvVar(); err != nil {
	// 	t.Errorf("Failed to load environment: %v", err)
	// 	return
	// }
	// setEnvVarForTesting(t)

	// if funcUsingEnvVar() != "dev" {
	// 	t.Errorf("Failed to load environment: expected dev, got %v", funcUsingEnvVar())
	// 	return
	// }

	var sms models.SmsModel // a == Student{"", 0}
	sms.Message = "12345"
	sms.ReciverPhoneNumber = "0587198145"
	sms.SenderName = "e-vrit"
	sms.SendingType = 1

	var isms smsmessage.ISms

	var flashySms flashy.FlashySmsModel
	flashySms.FlashyUrl = "https://api.flashy.app/"
	flashySms.SmsFlashyUrl = "messages/sms"
	flashySms.ContentType = "application/json"
	flashySms.Key = "vBbmiffyB4kaIrN2zCfa4luJe4Bbbmw7"
	//for getting error
	// flashySms.Key = "12345"
	flashySms.Data = []byte(`{"from": "Test User", "to": "test@example.com", "message": "test@example.com"}`)

	isms = &flashySms

	res := isms.SendSms(sms)
	if res == false {
		t.Error("test error")
	}
	fmt.Println(isms)
}
