package smsmessage

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	var isms ISms
	var err, errSms error

	if sms.SendingType > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "sending type is not exists"})
		return
	}

	token := os.Getenv("TOKEN")
	tokenReq := ctx.Request.Header["Token"][0]
	if bcrypt.CompareHashAndPassword([]byte(tokenReq), []byte(token)) != nil || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token not valid"})
		return
	}

	switch sms.SendingType {
	case enums.Flashy:
		isms, err = flashy.NewFlashySmsModel()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, errSms = isms.SendSms(sms)

	case enums.Inforu:
		isms, err = inforu.NewInforuSmsModel()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, errSms = isms.SendSms(sms)
	}

	if errSms != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errSms.Error()})
		return
	}
}
