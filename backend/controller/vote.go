package controller

import (
	"backend/dao/redis"
	"backend/logic"
	"backend/models"
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type VoteData struct {
	// UserID int Get the current user's ID from the request
	PostID    string `json:"post_id,string"`   // Post ID
	Direction int    `json:"direction,string"` // Vote direction (1 for upvote, -1 for downvote, 0 for cancel vote)
}

// UnmarshalJSON implements a custom UnmarshalJSON method for VoteData type
func (v *VoteData) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		PostID    string `json:"post_id"`
		Direction int    `json:"direction"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.PostID) == 0 {
		err = errors.New("Missing required field post_id")
	} else if required.Direction == 0 {
		err = errors.New("Missing required field direction")
	} else {
		v.PostID = required.PostID
		v.Direction = required.Direction
	}
	return
}

// VoteHandler handles the voting process
func VoteHandler(c *gin.Context) {
	// Parameter validation: which post to vote for and the vote direction
	vote := new(models.VoteDataForm)
	if err := c.ShouldBindJSON(&vote); err != nil {
		errs, ok := err.(validator.ValidationErrors) // Type assertion
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) // Translate and remove the struct prefix in error messages
		ResponseErrorWithMsg(c, CodeInvalidParams, errData)
		return
	}
	// Get the ID of the current user making the request
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	// Business logic for voting
	if err := logic.VoteForPost(userID, vote); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		switch err {
		case redis.ErrVoteRepeated: // Duplicate vote
			ResponseError(c, ErrVoteRepeated)
		case redis.ErrorVoteTimeExpire: // Vote time expired
			ResponseError(c, ErrorVoteTimeExpire)
		default:
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	ResponseSuccess(c, nil)
}
