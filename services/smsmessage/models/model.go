package models



type TemplateParameter struct {

	Name  string `json:"Name"`

	Type  string `json:"Type"`

	Value string `json:"Value"`

}



type SmsModel struct {

	SendingType        int      `json:"SendingType"`

	SenderName         string   `json:"SenderName"`

	Message            string   `json:"Message"`

	ReciverPhoneNumber string   `json:"ReciverPhoneNumber"`

	Token              string   `json:"Token"`

	Channels           []string `json:"Channels,omitempty"`



	// WhatsApp-specific fields

	TemplateId            int                 `json:"TemplateId,omitempty"`

	TemplateParameters    []TemplateParameter `json:"TemplateParameters,omitempty"`

	RecipientCustomFields map[string]string   `json:"RecipientCustomFields,omitempty"`

}



type ChannelResult struct {

	Channel string `json:"channel"`

	Success bool   `json:"success"`

	Error   string `json:"error,omitempty"`

}

