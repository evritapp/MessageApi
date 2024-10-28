package models

type SmsModel struct {
	SendingType        int
	SenderName         string
	Message            string
	ReciverPhoneNumber string
}