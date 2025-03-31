package models

import (
	"encoding/json"
	"errors"
	"time"
)

// Post Structure for Post. Concept of memory alignment: fields of the same type are aligned to reduce the memory size occupied by variables.
type Post struct {
	PostID      uint64    `json:"post_id,string" db:"post_id"`
	AuthorId    uint64    `json:"author_id" db:"author_id"`
	CommunityID uint64    `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"-" db:"create_time"`
	UpdateTime  time.Time `json:"-" db:"update_time"`
}

// UnmarshalJSON Custom UnmarshalJSON method for the Post type
func (p *Post) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Title       string `json:"title" db:"title"`
		Content     string `json:"content" db:"content"`
		CommunityID int64  `json:"community_id" db:"community_id"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.Title) == 0 {
		err = errors.New("Post title cannot be empty")
	} else if len(required.Content) == 0 {
		err = errors.New("Post content cannot be empty")
	} else if required.CommunityID == 0 {
		err = errors.New("Community not specified")
	} else {
		p.Title = required.Title
		p.Content = required.Content
		p.CommunityID = uint64(required.CommunityID)
	}
	return
}

// ApiPostDetail Structure for the detailed response of a post
type ApiPostDetail struct {
	*Post                                  // Embedded Post structure
	*CommunityDetailRes `json:"community"` // Embedded community information
	AuthorName          string             `json:"author_name"`
	VoteNum             int64              `json:"vote_num"` // Number of votes
	//CommunityName string `json:"community_name"`
}

type Page struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
}

type ApiPostDetailRes struct {
	Page Page             `json:"page"`
	List []*ApiPostDetail `json:"list"`
}
