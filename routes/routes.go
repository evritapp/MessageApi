package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"messageapi.e-vrit.co.il/services/smsmessage"
)

func SmsRoutes(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.POST("/sms/", smsmessage.SendSms)
	router.POST("/test/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})
}
