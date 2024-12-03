package inforu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func (inforuSmsModel InforuSmsModel) SendSms(sms models.SmsModel) (bool, error) {

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

	var jsonData bytes.Buffer
	if err := json.NewEncoder(&jsonData).Encode(data); err != nil {
		fmt.Printf("could not encode json: %s\n", err)
		return false, err
	}
	//Use Marshal
	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Printf("could not marshal json: %s\n", err)
	// 	return false
	// }
	// req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewReader(jsonData))

	fullUrl := inforuSmsModel.InforuUrl + inforuSmsModel.SmsInforuUrl
	req, err := http.NewRequest(http.MethodPost, fullUrl, &jsonData)

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
		return false, err
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
		return false, err
	}
	defer resp.Body.Close()
	//Use UnMarshal
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Printf("client: could not read response body: %s\n", err)
	// 	return false
	// }
	// var response map[string]interface{}
	// if err := json.Unmarshal(body, &response); err != nil {
	// 	fmt.Printf("client: could not parse response: %s\n", err)
	// 	return false
	// }

	decoder := json.NewDecoder(resp.Body)
	var response map[string]interface{}
	if err := decoder.Decode(&response); err != nil {
		fmt.Printf("client: could not parse response: %s\n", err)
		return false, err
	}

	if success, ok := response["StatusId"].(float64); ok && success == 1 {
		fmt.Println("Message sent successfully!")
		return true, nil
	}

	if errorRes, ok := response["StatusId"].(float64); ok && errorRes == -1 {
		fmt.Println("Failed to send message. Errors:")
		return false, err
	}

	return false, errors.New("unknown exception")
}
