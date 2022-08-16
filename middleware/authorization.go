package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorize(c *gin.Context) {
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		if auth, ok := c.GetQuery("authorization"); ok {
			jwtToken = auth
		}
	}

	if jwtToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "Plz login",
		})
		return
	}

	fmt.Printf("User Authorized! %v\n", jwtToken)
}
