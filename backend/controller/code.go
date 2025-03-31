package controller

type MyCode int64

const (
	CodeSuccess           MyCode = 1000
	CodeInvalidParams     MyCode = 1001
	CodeUserExist         MyCode = 1002
	CodeUserNotExist      MyCode = 1003
	CodeInvalidPassword   MyCode = 1004
	CodeServerBusy        MyCode = 1005
	CodeInvalidToken      MyCode = 1006
	CodeInvalidAuthFormat MyCode = 1007
	CodeNotLogin          MyCode = 1008
	ErrVoteRepeated       MyCode = 1009
	ErrorVoteTimeExpire   MyCode = 1010
)

var msgFlags = map[MyCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "Invalid request parameters",
	CodeUserExist:       "Username already exists",
	CodeUserNotExist:    "User does not exist",
	CodeInvalidPassword: "Invalid username or password",
	CodeServerBusy:      "Server is busy",

	CodeInvalidToken:      "Invalid token",
	CodeInvalidAuthFormat: "Incorrect authentication format",
	CodeNotLogin:          "Not logged in",
	ErrVoteRepeated:       "Do not vote repeatedly",
	ErrorVoteTimeExpire:   "Voting time has expired",
}

func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
