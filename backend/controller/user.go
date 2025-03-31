package controller

import (
	"backend/dao/mysql"
	"backend/logic"
	"backend/models"
	"backend/pkg/jwt"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignUpHandler handles user registration
func SignUpHandler(c *gin.Context) {
	// 1. Get request parameters
	var fo *models.RegisterForm
	// 2. Validate data
	if err := c.ShouldBindJSON(&fo); err != nil {
		// Invalid request parameters, return response directly
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// Check if err is of type validator.ValidationErrors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// Return if it's not of type validator.ValidationErrors
			ResponseError(c, CodeInvalidParams) // Invalid request parameters
			return
		}
		// Translate the validator.ValidationErrors type errors
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return // Error translation
	}
	fmt.Printf("fo: %v\n", fo)
	// 3. Business logic - Register user
	if err := logic.SignUp(fo); err != nil {
		zap.L().Error("logic.signup failed", zap.Error(err))
		if err.Error() == mysql.ErrorUserExit {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// Return response
	ResponseSuccess(c, nil)
}

// LoginHandler handles user login
func LoginHandler(c *gin.Context) {
	// 1. Get request parameters and validate
	var u *models.LoginForm
	if err := c.ShouldBindJSON(&u); err != nil {
		// Invalid request parameters, return response directly
		zap.L().Error("Login with invalid param", zap.Error(err))
		// Check if err is of type validator.ValidationErrors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// Return if it's not of type validator.ValidationErrors
			ResponseError(c, CodeInvalidParams) // Invalid request parameters
			return
		}
		// Translate the validator.ValidationErrors type errors
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2. Business logic - Login
	user, err := logic.Login(u)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", u.UserName), zap.Error(err))
		if err.Error() == mysql.ErrorUserNotExit {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 3. Return response
	ResponseSuccess(c, gin.H{
		"user_id":       fmt.Sprintf("%d", user.UserID), // The maximum value recognized by JS: id value greater than 1<<53-1  int64: i<<63-1
		"user_name":     user.UserName,
		"access_token":  user.AccessToken,
		"refresh_token": user.RefreshToken,
	})
}

// RefreshTokenHandler refreshes accessToken
func RefreshTokenHandler(c *gin.Context) {
	rt := c.Query("refresh_token")
	// There are three ways the client can carry the Token: 1. In the request header 2. In the request body 3. In the URI
	// Here, we assume the Token is in the Authorization header, prefixed with Bearer
	// The specific implementation should be based on your business logic
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithMsg(c, CodeInvalidToken, "Missing Auth Token in request header")
		c.Abort()
		return
	}
	// Split by space
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithMsg(c, CodeInvalidToken, "Invalid token format")
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	zap.L().Error("jwt.RefreshToken failed", zap.Error(err))
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}
