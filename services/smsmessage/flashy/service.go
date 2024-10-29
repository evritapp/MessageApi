package flashy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"messageapi.e-vrit.co.il/services/smsmessage/models"
)

func NewFlashySmsModel() (*FlashySmsModel, error) {
	fmt.Println("NewFlashySmsModel")
	err := godotenv.Load(".env")
	if err != nil {
		// log.Fatalf("Error loading .env file")
		return nil, err
	}

	return &FlashySmsModel{
		FlashyUrl:    os.Getenv("FLASHY_URL"),
		SmsFlashyUrl: "messages/sms",
		ContentType:  "application/json",
		Key:          os.Getenv("FLASHY_KEY"),
		Data:         []byte(`{"from": "Test User", "to": "test@example.com", "message": "test@example.com"}`),
	}, nil
}

func (flashySmsModel FlashySmsModel) SendSms(sms models.SmsModel) bool {

	// dataStr:= fmt.Sprintf(`"message": {{"from": "%v", "to": "%v", "message": "%v"}}`, sms.SenderName, sms.ReciverPhoneNumber, sms.Message)
	// data := []byte(dataStr)
	data := map[string]interface{}{
		"message": map[string]interface{}{
			"from":    sms.SenderName,
			"to":      sms.ReciverPhoneNumber,
			"message": sms.Message,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return false
	}
	// strJsonData = fmt.Sprintf("%v", jsonData)
	// client := &http.Client{}
	// req, err := http.NewRequest(http.MethodPost, flashySmsModel.FlashyUrl, bytes.NewBuffer(data))
	fullUrl := flashySmsModel.FlashyUrl + flashySmsModel.SmsFlashyUrl
	req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewReader(jsonData))

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
		return false
	}

	req.Header.Add("Content-Type", flashySmsModel.ContentType)
	req.Header.Add("x-api-key", flashySmsModel.Key)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	// res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("client: could not parse response: %s\n", err)
		return false
	}

	if success, ok := response["success"].(bool); ok && success {
		fmt.Println("Message sent successfully!")
		return true
	}

	if errors, ok := response["errors"].(map[string]interface{}); ok {
		fmt.Println("Failed to send message. Errors:")
		for key, val := range errors {
			fmt.Printf("%s: %v\n", key, val)
		}
	}

	return false
}
