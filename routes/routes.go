package routes

import (
	"github.com/gin-gonic/gin"
	"messageapi.e-vrit.co.il/services/smsmessage"
)

func SmsRoutes(router *gin.Engine) {

	router.POST("/sms/", smsmessage.SendSms)
}
