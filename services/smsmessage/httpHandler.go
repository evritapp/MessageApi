package smsmessage

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"messageapi.e-vrit.co.il/enums"
	"messageapi.e-vrit.co.il/services/smsmessage/flashy"
	"messageapi.e-vrit.co.il/services/smsmessage/inforu"
	"messageapi.e-vrit.co.il/services/smsmessage/models"
)

func SendSms(ctx *gin.Context) {

	var sms models.SmsModel
	if err := ctx.ShouldBindJSON(&sms); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := os.Getenv("TOKEN")
	tokenReq := ctx.GetHeader("Token")
	if token == "" || tokenReq == "" || token != tokenReq {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token not valid"})
		return
	}

	channels := sms.Channels
	if len(channels) == 0 {
		channels = []string{"sms"}
	}

	var results []models.ChannelResult
	hasSuccess := false

	for _, channel := range channels {
		var result models.ChannelResult

		switch channel {
		case "sms":
			result = sendSmsChannel(sms)
		case "whatsapp":
			result = sendWhatsAppChannel(sms)
		default:
			result = models.ChannelResult{
				Channel: channel,
				Error:   "invalid channel, use 'sms' or 'whatsapp'",
			}
		}

		if result.Success {
			hasSuccess = true
		}
		results = append(results, result)
	}

	if hasSuccess {
		ctx.JSON(http.StatusOK, gin.H{"results": results})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"results": results})
	}
}

func sendSmsChannel(sms models.SmsModel) models.ChannelResult {
	result := models.ChannelResult{Channel: "sms"}

	var isms ISms
	var err error

	switch sms.SendingType {
	case enums.Inforu:
		isms, err = inforu.NewInforuModel()
	case enums.Flashy:
		isms, err = flashy.NewFlashySmsModel()
	default:
		result.Error = "sending type does not exist"
		return result
	}

	if err != nil {
		result.Error = err.Error()
		return result
	}

	_, errSms := isms.SendSms(sms)
	if errSms != nil {
		result.Error = errSms.Error()
		return result
	}

	result.Success = true
	return result
}

func sendWhatsAppChannel(sms models.SmsModel) models.ChannelResult {
	result := models.ChannelResult{Channel: "whatsapp"}

	if sms.SendingType != enums.Inforu {
		result.Error = "WhatsApp is only supported via InfoRU"
		return result
	}

	if sms.TemplateId == 0 {
		result.Error = "TemplateId is required for WhatsApp"
		return result
	}

	isms, err := inforu.NewInforuModel()
	if err != nil {
		result.Error = err.Error()
		return result
	}

	_, errSms := isms.SendWhatsApp(sms)
	if errSms != nil {
		result.Error = errSms.Error()
		return result
	}

	result.Success = true
	return result
}
