package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondOk() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Good probe",
			"data":    gin.H{},
		})
	}
}
