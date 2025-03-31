package mysql

import (
	"backend/models"
	"backend/pkg/snowflake"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// Encapsulate each database operation as a function
// Logic layer will call these functions based on business requirements

const secret = "huchao.vip"

// encryptPassword Encrypt the password
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}

// CheckUserExist Check if a user with the given username exists
func CheckUserExist(username string) (error error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New(ErrorUserExit) // User exists
	}
	return
}

// InsertUser Register a new user by inserting into the database
func InsertUser(user models.User) (error error) {
	// Encrypt the password
	user.Password = encryptPassword([]byte(user.Password))
	// Execute the SQL to insert the user into the database
	sqlstr := `insert into user(user_id,username,password,email,gender) values(?,?,?,?,?)`
	_, err := db.Exec(sqlstr, user.UserID, user.UserName, user.Password, user.Email, user.Gender)
	return err
}

func Register(user *models.User) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int64
	err = db.Get(&count, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count > 0 {
		// User already exists
		return errors.New(ErrorUserExit)
	}
	// Generate user_id
	userID, err := snowflake.GetID()
	if err != nil {
		return ErrorGenIDFailed
	}
	// Generate encrypted password
	password := encryptPassword([]byte(user.Password))
	// Insert the user into the database
	sqlStr = "insert into user(user_id, username, password) values (?,?,?)"
	_, err = db.Exec(sqlStr, userID, user.UserName, password)
	return
}

// Login Login business logic
func Login(user *models.User) (err error) {
	originPassword := user.Password // Keep track of the original password (the password the user logs in with)
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.UserName)
	// Error while querying the database
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	// User does not exist
	if err == sql.ErrNoRows {
		return errors.New(ErrorUserNotExit)
	}
	// Compare the encrypted password with the password fetched from the database
	password := encryptPassword([]byte(originPassword))
	if user.Password != password {
		return errors.New(ErrorPasswordWrong)
	}
	return nil
}

// GetUserByID Query author information by ID
func GetUserByID(id uint64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}
