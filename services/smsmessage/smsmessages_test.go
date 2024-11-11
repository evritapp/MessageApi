package smsmessage_test

import (
	"fmt"
	"os"
	"path/filepath"

	"testing"

	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"
	"messageapi.e-vrit.co.il/services/smsmessage"
	"messageapi.e-vrit.co.il/services/smsmessage/flashy"
	"messageapi.e-vrit.co.il/services/smsmessage/inforu"
	"messageapi.e-vrit.co.il/services/smsmessage/models"
)


var testFlashyCases = []struct {
	name          string
	expected      bool
	expectedToken bool
	smsModel      models.SmsModel
	FlashySms     flashy.FlashySmsModel
	Token         string
}{
	{
		name:          "flashy sms- success",
		expected:      true,
		expectedToken: true,
		smsModel: models.SmsModel{
			Message:            " הצלחה פלאשי",
			ReciverPhoneNumber: "0587198145",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		FlashySms: flashy.FlashySmsModel{
			FlashyUrl:    "https://api.flashy.app/",
			SmsFlashyUrl: "messages/sms",
			ContentType:  "application/json",
			Key:          "vBbmiffyB4kaIrN2zCfa4luJe4Bbbmw7",
		},
		Token: "1q2w3e4r",
	},

	{
		name:          "flashy sms- not success key not valid",
		expected:      false,
		expectedToken: true,
		smsModel: models.SmsModel{
			Message:            " כשלון פלאשי",
			ReciverPhoneNumber: "0587198145",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		FlashySms: flashy.FlashySmsModel{
			FlashyUrl:    "https://api.flashy.app/",
			SmsFlashyUrl: "messages/sms",
			ContentType:  "application/json",
			Key:          "vBbmiffyB4kaIrN2zC555fa4khjkhkhkluJe4Bbbmw7",
		},
		Token: "1q2w3e4r",
	},
	{
		name:          "flashy sms- not success- token not valid",
		expected:      true,
		expectedToken: false,
		smsModel: models.SmsModel{
			Message:            " הצלחה פלאשי",
			ReciverPhoneNumber: "0587198145",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		FlashySms: flashy.FlashySmsModel{
			FlashyUrl:    "https://api.flashy.app/",
			SmsFlashyUrl: "messages/sms",
			ContentType:  "application/json",
			Key:          "vBbmiffyB4kaIrN2zCfa4luJe4Bbbmw7",
		},
		Token: "123456",
	},
}

func TestFlashySendSms(t *testing.T) {

	var isms smsmessage.ISms
	envToken := GetToken()
	for _, tc := range testFlashyCases {

		asserts := assert.New(t)

		isToken := envToken == tc.Token

		asserts.Equal(tc.expectedToken, isToken)

		isms = &tc.FlashySms

		res, _ := isms.SendSms(tc.smsModel)
		asserts.Equal(tc.expected, res)

	}

}

var testInforuCases = []struct {
	name          string
	expected      bool
	expectedToken bool
	smsModel      models.SmsModel
	InforuSms     inforu.InforuSmsModel
	Token         string
}{
	{
		name:          "inforu sms- success",
		expected:      true,
		expectedToken: true,
		smsModel: models.SmsModel{
			Message:            " הצלחה אינפוריו",
			ReciverPhoneNumber: "0587198145",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		InforuSms: inforu.InforuSmsModel{
			InforuUrl:     "https://capi.inforu.co.il/api/v2/",
			SmsInforuUrl:  "SMS/SendSms",
			ContentType:   "application/json",
			Authorization: "Basic WWFxdWllbDoxZDk3Zjc4Yi1jNTIzLTRjMDctOTU5Ni1jNjk4YzdiMzQ2YzU=",
		},
		Token: "1q2w3e4r",
	},

	{
		name:          "inforu sms- not success Authorization not valid",
		expected:      false,
		expectedToken: true,
		smsModel: models.SmsModel{
			Message:            " כשלון אינפוריו",
			ReciverPhoneNumber: "0587198145",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		InforuSms: inforu.InforuSmsModel{
			InforuUrl:     "https://capi.inforu.co.il/api/v2/",
			SmsInforuUrl:  "SMS/SendSms",
			ContentType:   "application/json",
			Authorization: "Basic WWFxdWllbDoxZDk3Zjc4Yi1jdfgdgNTIzLTRjMDctOTU5Nli1jNjk4YzdiMzQ2YzU=",
		},
		Token: "1q2w3e4r",
	},
	{
		name:          "inforu sms- not success token not valid",
		expected:      true,
		expectedToken: false,
		smsModel: models.SmsModel{
			Message:            " הצלחה אינפוריו",
			ReciverPhoneNumber: "0587198145",
			SenderName:         "e-vrit",
			SendingType:        1,
		},
		InforuSms: inforu.InforuSmsModel{
			InforuUrl:     "https://capi.inforu.co.il/api/v2/",
			SmsInforuUrl:  "SMS/SendSms",
			ContentType:   "application/json",
			Authorization: "Basic WWFxdWllbDoxZDk3Zjc4Yi1jNTIzLTRjMDctOTU5Ni1jNjk4YzdiMzQ2YzU=",
		},
		Token: "21245454",
	},
}

func TestInforuSendSms(t *testing.T) {

	var isms smsmessage.ISms
	envToken := GetToken()
	for _, tc := range testInforuCases {

		asserts := assert.New(t)
		isToken := envToken == tc.Token

		asserts.Equal(tc.expectedToken, isToken)

		isms = &tc.InforuSms
		res, _ := isms.SendSms(tc.smsModel)
		asserts.Equal(tc.expected, res)

	}
}

func GetToken() string {

	wd, err := filepath.Abs(".")
	if err != nil {
		fmt.Println("Error:", err)
	}
	filepath.Abs(".")
	envPath := filepath.Join(filepath.Dir(filepath.Dir(wd)), ".env")
	if err := godotenv.Load(envPath); err != nil {
		fmt.Println("Error loading .env file:", err)
		return ""
	}

	t := os.Getenv("TOKEN")
	return t
}
