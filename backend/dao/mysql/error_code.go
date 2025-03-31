package mysql

import "errors"

// Define business status
var (
	ErrorUserExit      = "User already exists"
	ErrorUserNotExit   = "User does not exist"
	ErrorPasswordWrong = "Incorrect password"
	ErrorGenIDFailed   = errors.New("Failed to generate user ID")
	ErrorInvalidID     = "Invalid ID"
	ErrorQueryFailed   = "Failed to query data"
	ErrorInsertFailed  = errors.New("Failed to insert data")
)
