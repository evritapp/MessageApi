package models

type TemplateParameter struct {
	Name string `json:"Name"`

	Type string `json:"Type"`

	Value string `json:"Value"`
}

type Button struct {
	ButtonIndex int    `json:"ButtonIndex"`
	FieldSource string `json:"FieldSource"`
	FieldName   string `json:"FieldName"`
	Value       string `json:"Value,omitempty"`
	Payload     string `json:"Payload,omitempty"`
	PayloadAlt  string `json:"payload,omitempty"`
}

type SmsModel struct {
	SendingType int `json:"SendingType"`

	SenderName string `json:"SenderName"`

	Message string `json:"Message"`

	ReciverPhoneNumber string `json:"ReciverPhoneNumber"`

	Token string `json:"Token"`

	Channels []string `json:"Channels,omitempty"`

	// WhatsApp-specific fields

	TemplateId int `json:"TemplateId,omitempty"`

	TemplateParameters []TemplateParameter `json:"TemplateParameters,omitempty"`

	RecipientCustomFields map[string]string `json:"RecipientCustomFields,omitempty"`
	Buttons               []Button          `json:"Buttons,omitempty"`
}

type ChannelResult struct {
	Channel string `json:"channel"`

	Success bool `json:"success"`

	Error string `json:"error,omitempty"`
}
