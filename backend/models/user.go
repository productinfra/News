package models

import (
	"encoding/json"
	"errors"
)

// User Defines the structure of request parameters for the user
type User struct {
	UserID       uint64 `json:"user_id,string" db:"user_id"` // Specifies the use of lowercase user_id in JSON serialization/deserialization
	UserName     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	Email        string `json:"email" db:"gender"`  // Email
	Gender       int    `json:"gender" db:"gender"` // Gender
	AccessToken  string
	RefreshToken string
}

// UnmarshalJSON Custom UnmarshalJSON method for the User type
func (u *User) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		UserName string `json:"username" db:"username"`
		Password string `json:"password" db:"password"`
		Email    string `json:"email" db:"gender"`  // Email
		Gender   int    `json:"gender" db:"gender"` // Gender
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.UserName) == 0 {
		err = errors.New("Missing required field: username")
	} else if len(required.Password) == 0 {
		err = errors.New("Missing required field: password")
	} else {
		u.UserName = required.UserName
		u.Password = required.Password
		u.Email = required.Email
		u.Gender = required.Gender
	}
	return
}

// RegisterForm Structure for registration request parameters
type RegisterForm struct {
	UserName        string `json:"username" binding:"required"`  // Username
	Email           string `json:"email" binding:"required"`     // Email
	Gender          int    `json:"gender" binding:"oneof=0 1 2"` // Gender 0: Unknown, 1: Male, 2: Female
	Password        string `json:"password" binding:"required"`  // Password
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// LoginForm Structure for login request parameters
type LoginForm struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UnmarshalJSON Custom UnmarshalJSON method for the RegisterForm type
func (r *RegisterForm) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		UserName        string `json:"username"`
		Email           string `json:"email"`    // Email
		Gender          int    `json:"gender"`   // Gender 0: Unknown, 1: Male, 2: Female
		Password        string `json:"password"` // Password
		ConfirmPassword string `json:"confirm_password"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.UserName) == 0 {
		err = errors.New("Missing required field: username")
	} else if len(required.Password) == 0 {
		err = errors.New("Missing required field: password")
	} else if len(required.Email) == 0 {
		err = errors.New("Missing required field: email")
	} else if required.Password != required.ConfirmPassword {
		err = errors.New("Passwords do not match")
	} else {
		r.UserName = required.UserName
		r.Email = required.Email
		r.Gender = required.Gender
		r.Password = required.Password
		r.ConfirmPassword = required.ConfirmPassword
	}
	return
}

// VoteDataForm Structure for vote data
type VoteDataForm struct {
	// UserID int To get the current user's ID from the request
	PostID    string `json:"post_id" binding:"required"`              // Post ID
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // Vote direction (1: Agree, -1: Disagree, 0: Cancel vote)
}

// UnmarshalJSON Custom UnmarshalJSON method for the VoteDataForm type
func (v *VoteDataForm) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		PostID    string `json:"post_id"`
		Direction int8   `json:"direction"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.PostID) == 0 {
		err = errors.New("Missing required field: post_id")
	} else if required.Direction == 0 {
		err = errors.New("Missing required field: direction")
	} else {
		v.PostID = required.PostID
		v.Direction = required.Direction
	}
	return
}
