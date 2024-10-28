package smsmessage

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"messageapi.e-vrit.co.il/enums"
	"messageapi.e-vrit.co.il/services/smsmessage/flashy"
	"messageapi.e-vrit.co.il/services/smsmessage/inforu"
	"messageapi.e-vrit.co.il/services/smsmessage/models"
)

func SendSms(ctx *gin.Context) {
	fmt.Println("controller", ctx)
	var sms models.SmsModel
	// sms.SendingType = 1
	if err := ctx.ShouldBindJSON(&sms); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var isms ISms
	var err error
	switch sms.SendingType {
	case enums.Flashy:
		isms, err = flashy.NewFlashySmsModel()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		isms.SendSms(sms)
	case enums.Inforu:
		isms, err = inforu.NewInforuSmsModel()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		isms.SendSms(sms)
	}

}
