package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n0l3r/limitron"
)

func RateLimitMiddleware(limiter limitron.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		allowed, err := limiter.Allow(key)
		if err != nil || !allowed {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
