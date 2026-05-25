package inforu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"messageapi.e-vrit.co.il/services/smsmessage/models"
)

func NewInforuModel() (*InforuModel, error) {
	url := os.Getenv("INFORU_URL")
	auth := os.Getenv("INFORU_AUT")
	if url == "" || auth == "" {
		return nil, errors.New("inforu: INFORU_URL and INFORU_AUT must be set")
	}

	return &InforuModel{
		InforuUrl:     url,
		ContentType:   "application/json",
		Authorization: auth,
	}, nil
}

func (m InforuModel) SendSms(sms models.SmsModel) (bool, error) {
	data := map[string]interface{}{
		"Message": sms.Message,
		"Recipients": []map[string]interface{}{
			{"Phone": sms.ReciverPhoneNumber},
		},
		"Settings": map[string]interface{}{
			"Sender": sms.SenderName,
		},
	}

	return m.sendRequest("SMS/SendSms", data)
}

func (m InforuModel) SendWhatsApp(sms models.SmsModel) (bool, error) {
	recipient := map[string]interface{}{
		"Phone": sms.ReciverPhoneNumber,
	}
	for k, v := range sms.RecipientCustomFields {
		recipient[k] = v
	}

	requestData := map[string]interface{}{
		"TemplateId": sms.TemplateId,
		"Recipients": []map[string]interface{}{recipient},
	}

	if sms.TemplateParameters != nil {
		requestData["TemplateParameters"] = sms.TemplateParameters
	} else {
		requestData["TemplateParameters"] = []models.TemplateParameter{}
	}

	if len(sms.Buttons) > 0 {
		buttons := make([]map[string]interface{}, 0, len(sms.Buttons))
		for _, b := range sms.Buttons {
			value := b.Value
			if value == "" {
				value = b.Payload
			}
			if value == "" {
				value = b.PayloadAlt
			}

			buttons = append(buttons, map[string]interface{}{
				"ButtonIndex": b.ButtonIndex,
				"FieldSource": b.FieldSource,
				"FieldName":   b.FieldName,
				"Value":       value,
			})
		}
		requestData["Buttons"] = buttons
	}

	data := map[string]interface{}{
		"Data": requestData,
	}

	return m.sendRequest("WhatsApp/SendWhatsApp", data)
}

func (m InforuModel) sendRequest(pathUrl string, data map[string]interface{}) (bool, error) {
	var jsonData bytes.Buffer
	if err := json.NewEncoder(&jsonData).Encode(data); err != nil {
		fmt.Printf("could not encode json: %s\n", err)
		return false, err
	}

	fullUrl := m.InforuUrl + pathUrl
	req, err := http.NewRequest(http.MethodPost, fullUrl, &jsonData)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return false, err
	}

	req.Header.Add("Content-Type", m.ContentType)
	req.Header.Add("Authorization", m.Authorization)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return false, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var response map[string]interface{}
	if err := decoder.Decode(&response); err != nil {
		fmt.Printf("client: could not parse response: %s\n", err)
		return false, err
	}

	if statusId, ok := response["StatusId"].(float64); ok {
		if statusId == 1 {
			fmt.Println("Message sent successfully!")
			return true, nil
		}
		if statusId == -1 {
			fmt.Println("Failed to send message.")
			return false, errors.New("inforu returned error status")
		}
	}

	return false, errors.New("unknown exception")
}
