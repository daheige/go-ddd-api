package utils

import (
	"log"
	"runtime/debug"
)

// Recover catch server recover
func Recover() {
	if err := recover(); err != nil {
		log.Println("exec panic", map[string]interface{}{
			"error":       err,
			"error_trace": string(debug.Stack()),
		})
	}
}
