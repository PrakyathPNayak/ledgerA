package middleware

import (
	"fmt"
	"ledgerA/internal/dto"
	firebasepkg "ledgerA/pkg/firebase"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	contextFirebaseUIDKey = "firebase_uid"
	contextEmailKey       = "email"
)

// AuthMiddleware verifies Firebase JWT from Authorization header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			dto.Error(c, 401, "ERR_UNAUTHORIZED", "missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			dto.Error(c, 401, "ERR_UNAUTHORIZED", "invalid Authorization header format")
			c.Abort()
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			dto.Error(c, 401, "ERR_UNAUTHORIZED", "missing bearer token")
			c.Abort()
			return
		}

		verified, err := firebasepkg.VerifyIDToken(c.Request.Context(), token)
		if err != nil {
			dto.Error(c, 401, "ERR_UNAUTHORIZED", "token verification failed")
			c.Abort()
			return
		}

		c.Set(contextFirebaseUIDKey, verified.UID)
		email := ""
		if value, ok := verified.Claims["email"]; ok {
			email = fmt.Sprintf("%v", value)
		}
		c.Set(contextEmailKey, email)
		c.Next()
	}
}

// GetFirebaseUID returns firebase uid from context.
func GetFirebaseUID(c *gin.Context) (string, error) {
	value, ok := c.Get(contextFirebaseUIDKey)
	if !ok {
		return "", fmt.Errorf("middleware.GetFirebaseUID: firebase uid not found in context")
	}
	uid, ok := value.(string)
	if !ok || uid == "" {
		return "", fmt.Errorf("middleware.GetFirebaseUID: invalid firebase uid in context")
	}
	return uid, nil
}
