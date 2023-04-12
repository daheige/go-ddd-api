package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error is wrapped Respond when error response
func Error(c *gin.Context, code int, message string, err ...error) {
	if code <= 0 {
		code = http.StatusInternalServerError
	}

	m := map[string]interface{}{
		"code":    code,
		"message": message,
	}

	if len(err) > 0 && err[0] != nil {
		m["trace_error"] = err[0].Error()
	}

	c.JSON(http.StatusOK, m)
}

// Success http api success
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": message,
		"data":    data,
	})
}

// JSON is wrapped Respond when success response
func JSON(c *gin.Context, code int, data interface{}, msg ...string) {
	var message string
	if len(msg) > 0 {
		message = msg[0]
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	})
}
