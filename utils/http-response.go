package utils

import (
	"github.com/gin-gonic/gin"
)

func SendSuccessResponse(con *gin.Context, statusCode int, result any) {
	con.JSON(statusCode, result)
}

func SendErrorResponse(con *gin.Context, statusCode int, err error) {
	con.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}
