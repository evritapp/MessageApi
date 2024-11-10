package smsmessage_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"messageapi.e-vrit.co.il/services/smsmessage"
	"messageapi.e-vrit.co.il/services/smsmessage/flashy"
	"messageapi.e-vrit.co.il/services/smsmessage/inforu"
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

var testFlashyCases = []struct {
	name      string
	expected  bool
	smsModel  models.SmsModel
	FlashySms flashy.FlashySmsModel
}{
	{
		name:     "flashy sms- success",
		expected: true,
		smsModel: models.SmsModel{
			Message:            " הצלחה פלאשי",
			ReciverPhoneNumber: "0526012123",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		FlashySms: flashy.FlashySmsModel{
			FlashyUrl:    "https://api.flashy.app/",
			SmsFlashyUrl: "messages/sms",
			ContentType:  "application/json",
			Key:          "vBbmiffyB4kaIrN2zCfa4luJe4Bbbmw7",
		},
	},

	{
		name:     "flashy sms- not success",
		expected: false,
		smsModel: models.SmsModel{
			Message:            " כשלון פלאשי",
			ReciverPhoneNumber: "0526012123",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		FlashySms: flashy.FlashySmsModel{
			FlashyUrl:    "https://api.flashy.app/",
			SmsFlashyUrl: "messages/sms",
			ContentType:  "application/json",
			Key:          "vBbmiffyB4kaIrN2zC555fa4khjkhkhkluJe4Bbbmw7",
		},
	},
}

func TestFlashySendSms(t *testing.T) {

	var isms smsmessage.ISms
	for _, tc := range testFlashyCases {

		asserts := assert.New(t)
		isms = &tc.FlashySms

		res := isms.SendSms(tc.smsModel)
		asserts.Equal(tc.expected, res)
	}

}

var testInforuCases = []struct {
	name      string
	expected  bool
	smsModel  models.SmsModel
	InforuSms inforu.InforuSmsModel
}{
	{
		name:     "inforu sms- success",
		expected: true,
		smsModel: models.SmsModel{
			Message:            "אינפוריו",
			ReciverPhoneNumber: "0526012123",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		InforuSms: inforu.InforuSmsModel{
			InforuUrl:     "https://capi.inforu.co.il/api/v2/",
			SmsInforuUrl:  "SMS/SendSms",
			ContentType:   "application/json",
			Authorization: "Basic WWFxdWllbDoxZDk3Zjc4Yi1jNTIzLTRjMDctOTU5Ni1jNjk4YzdiMzQ2YzU=",
		},
	},

	{
		name:     "inforu sms- not success",
		expected: true,
		smsModel: models.SmsModel{
			Message:            "אינפוריו",
			ReciverPhoneNumber: "0526012123",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		InforuSms: inforu.InforuSmsModel{
			InforuUrl:     "https://capi.inforu.co.il/api/v2/",
			SmsInforuUrl:  "SMS/SendSms",
			ContentType:   "application/json",
			Authorization: "Basic WWFxdWllbDoxZDk3Zjc4Yi1jdfgdgNTIzLTRjMDctOTU5Ni1jNjk4YzdiMzQ2YzU=",
		},
	},
}

func TestInforuSendSms(t *testing.T) {

	var isms smsmessage.ISms
	for _, tc := range testInforuCases {

		asserts := assert.New(t)
		isms = &tc.InforuSms

		res := isms.SendSms(tc.smsModel)
		asserts.Equal(tc.expected, res)
	}

}
