package redis

import "errors"

var (
	ErrorVoteTimeExpire = errors.New("Voting time has expired")
	ErrorVoted          = errors.New("You have already voted")
	ErrVoteRepeated     = errors.New("Repeated voting is not allowed")
)
