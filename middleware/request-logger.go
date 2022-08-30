package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger(param gin.LogFormatterParams) string {
	// your custom format
	return fmt.Sprintf("[%s] \"%s %s %d %s \"%s \"\n",
		param.TimeStamp.Format(time.RFC3339Nano),
		param.Method,
		param.Path,
		param.StatusCode,
		param.Latency,
		param.ErrorMessage,
	)
}
