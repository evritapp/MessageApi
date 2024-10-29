package inforu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"messageapi.e-vrit.co.il/services/smsmessage/models"
)

func NewInforuSmsModel() (*InforuSmsModel, error) {
	fmt.Println("NewInforuSmsModel")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		return nil, err
	}

	return &InforuSmsModel{
		InforuUrl:     os.Getenv("INFORU_URL"),
		SmsInforuUrl:  "SMS/SendSms",
		ContentType:   "application/json",
		Authorization: os.Getenv("INFORU_AUT"),
		Data:          []byte(""),
	}, nil
}

func (inforuSmsModel InforuSmsModel) SendSms(sms models.SmsModel) bool {

	data := map[string]interface{}{
		"Message": sms.Message,
		"Recipients": []map[string]interface{}{
			{
				"Phone": sms.ReciverPhoneNumber,
			},
		},
		"Settings": map[string]interface{}{
			"Sender": sms.SenderName,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return false
	}

	fullUrl := inforuSmsModel.InforuUrl + inforuSmsModel.SmsInforuUrl
	req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewReader(jsonData))

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
		return false
	}

	req.Header.Add("Content-Type", inforuSmsModel.ContentType)
	req.Header.Add("Authorization", inforuSmsModel.Authorization)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
		return false
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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
