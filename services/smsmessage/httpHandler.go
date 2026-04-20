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

	channel := ctx.DefaultQuery("channel", "sms")

	var isms ISms
	var err error

	switch channel {
	case "sms":
		switch sms.SendingType {
		case enums.Inforu:
			isms, err = inforu.NewInforuModel()
		case enums.Flashy:
			isms, err = flashy.NewFlashySmsModel()
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "sending type does not exist"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, errSms := isms.SendSms(sms)
		if errSms != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errSms.Error()})
			return
		}

	case "whatsapp":
		if sms.SendingType != enums.Inforu {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "WhatsApp is only supported via InfoRU"})
			return
		}
		if sms.TemplateId == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "TemplateId is required for WhatsApp"})
			return
		}
		isms, err = inforu.NewInforuModel()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, errSms := isms.SendWhatsApp(sms)
		if errSms != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errSms.Error()})
			return
		}

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel, use 'sms' or 'whatsapp'"})
		return
	}
}
