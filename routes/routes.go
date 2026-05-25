package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"messageapi.e-vrit.co.il/services/smsmessage"
)

func SmsRoutes(router *gin.Engine) {

	router.POST("/sms/", smsmessage.SendSms)
	router.POST("/test/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})
}
