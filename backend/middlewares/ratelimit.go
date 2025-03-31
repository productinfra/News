package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// The token bucket algorithm is similar to the leaky bucket principle. The token bucket adds tokens at a fixed rate, and as long as tokens are available in the bucket, the request can pass through. The token bucket supports rapid processing of burst traffic.
// For scenarios where tokens cannot be obtained from the bucket, we can choose to wait or reject the request and return an error.
// For implementing a token bucket in Go, you can refer to the github.com/juju/ratelimit library. This library supports multiple token bucket modes and is simple to use.

// For the registration location of this rate limiting middleware, we can register it at different places depending on the rate limiting strategy, such as:
// If rate limiting is required for the entire site, it can be registered as a global middleware.
// If rate limiting is needed for a specific group of routes, the middleware can be registered only for the corresponding route group.

// RateLimitMiddleware creates a token bucket with a specified fill rate and capacity
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// If no token is available, interrupt the request and return "rate limit..."
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		// If a token is obtained, proceed with the request
		c.Next()
	}
}
