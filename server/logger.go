package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Enumerated log levels (log levels match that of logrus)
const (
	Info = iota
	Warn
	Error
	Fatal
	Panic
)

// Log helper method to make it more convenient to create structured log messages
func LogRequest(level int, message string, latency time.Duration, status int, method string, route string) {
	logFields := log.Fields{
		"latency": latency,
		"status":  status,
		"method":  method,
		"route":   route,
	}
	switch level {
	case Info:
		log.WithFields(logFields).Info(message)
	case Warn:
		log.WithFields(logFields).Warn(message)
	case Error:
		log.WithFields(logFields).Error(message)
	case Fatal:
		log.WithFields(logFields).Fatal(message)
	case Panic:
		log.WithFields(logFields).Panic(message)
	}
	return
}

// Logger middleware - simply logs info about handled requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// before request

		c.Next()

		// after request
		latency := time.Since(t)

		// access the status we are sending
		status := c.Writer.Status()
		method := c.Request.Method
		route := c.Request.URL.String()

		var logLevel int
		var message string

		switch status {
		case http.StatusOK:
			fallthrough
		case http.StatusAccepted:
			logLevel = Info
			message = "Handled request successfully"
		case http.StatusNotFound:
			logLevel = Info
			message = "Unhandled route"
		case http.StatusInternalServerError:
			logLevel = Error
			message = "Internal server error"
		case http.StatusBadRequest:
			logLevel = Info
			message = "Bad request from client"
		default:
			logLevel = Warn
			message = "Encountered unexpected status"
		}
		LogRequest(logLevel, message, latency, status, method, route)

	}
}
