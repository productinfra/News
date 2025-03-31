package models

const (
	OrderTime  = "time"  // Order by time
	OrderScore = "score" // Order by score
)

// ParamPostList Parameters for getting a list of posts
type ParamPostList struct {
	Search      string `json:"search" form:"search"`               // Keyword search
	CommunityID uint64 `json:"community_id" form:"community_id"`   // Can be empty
	Page        int64  `json:"page" form:"page"`                   // Page number
	Size        int64  `json:"size" form:"size"`                   // Number of items per page
	Order       string `json:"order" form:"order" example:"score"` // Sorting criteria
}

// ParamGithubTrending Parameters for getting GitHub trending projects
type ParamGithubTrending struct {
	Language int   `json:"language" form:"language"` // Language: 0 - All, 1 - Go, 2 - Python, 3 - JavaScript, 4 - Java
	Page     int64 `json:"page" form:"page"`         // Page number
	Size     int64 `json:"size" form:"size"`         // Number of items per page
}
