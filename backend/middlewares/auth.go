package middlewares

import (
	"backend/controller"
	"backend/pkg/jwt"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware - JWT-based authentication middleware
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// There are three ways the client can send the token:
		// 1. In the request header
		// 2. In the request body
		// 3. In the URI
		// Here we assume the token is placed in the Authorization header, starting with Bearer
		// The specific implementation method depends on your actual business requirements
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseErrorWithMsg(c, controller.CodeInvalidToken, "Authorization Token is missing from the request header")
			c.Abort()
			return
		}
		// Split by space
		//&& parts[0] == "Bearer"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2) {
			controller.ResponseErrorWithMsg(c, controller.CodeInvalidToken, "Invalid token format")
			c.Abort()
			return
		}
		// parts[1] contains the tokenString, we use the previously defined function to parse the JWT
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			fmt.Println(err)
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// Save the userID of the current request to the context
		c.Set(controller.ContextUserIDKey, mc.UserID)
		c.Next() // The subsequent handlers can retrieve the current user's information using c.Get(ContextUserIDKey)
	}
}
