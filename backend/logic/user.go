package logic

import (
	"backend/dao/mysql"
	"backend/models"
	"backend/pkg/jwt"
	"backend/pkg/snowflake"
)

// SignUp registration business logic
func SignUp(p *models.RegisterForm) (error error) {
	// 1. Check if the user already exists
	err := mysql.CheckUserExist(p.UserName)
	if err != nil {
		// Error occurred during database query
		return err
	}

	// 2. Generate UID
	userId, err := snowflake.GetID()
	if err != nil {
		return mysql.ErrorGenIDFailed
	}
	// Construct a User instance
	u := models.User{
		UserID:   userId,
		UserName: p.UserName,
		Password: p.Password,
		Email:    p.Email,
		Gender:   p.Gender,
	}
	// 3. Save to database
	return mysql.InsertUser(u)
}

// Login login business logic
func Login(p *models.LoginForm) (user *models.User, error error) {
	user = &models.User{
		UserName: p.UserName,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// Generate JWT
	//return jwt.GenToken(user.UserID, user.UserName)
	accessToken, refreshToken, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return
	}
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	return
}
