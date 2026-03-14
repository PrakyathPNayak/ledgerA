package middleware

import (
	"fmt"
	"ledgerA/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RecoveryMiddleware captures panics and returns standardized 500 responses.
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Str("panic", fmt.Sprintf("%v", r)).
					Msg("panic recovered")

				dto.Error(c, 500, "ERR_INTERNAL", "internal server error")
				c.Abort()
			}
		}()

		c.Next()
	}
}
