package models

import "time"

// Community represents the Community structure
type Community struct {
	CommunityID   uint64 `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
}

// CommunityDetail represents the community details model
type CommunityDetail struct {
	CommunityID   uint64    `json:"community_id" db:"community_id"`
	CommunityName string    `json:"community_name" db:"community_name"`
	Introduction  string    `json:"introduction,omitempty" db:"introduction"` // omitempty means Introduction will not be displayed if it is empty
	CreateTime    time.Time `json:"create_time" db:"create_time"`
}

// CommunityDetailRes represents the response model for community details
type CommunityDetailRes struct {
	CommunityID   uint64 `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
	Introduction  string `json:"introduction,omitempty" db:"introduction"` // omitempty means Introduction will not be displayed if it is empty
	CreateTime    string `json:"create_time" db:"create_time"`
}
