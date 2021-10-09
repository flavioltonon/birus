package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "Request-ID"

// SetRequestID adds a request-id header to the request
func SetRequestID(ctx *gin.Context) {
	requestID := uuid.NewString()
	ctx.Request.Header.Set(RequestIDKey, requestID)
	ctx.Set(RequestIDKey, requestID)
	ctx.Next()
}
