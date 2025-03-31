package controller

import "backend/models"

// Since the response format of our API documentation is consistent, but the specific data type varies
type _ResponsePostList struct {
	Code    MyCode                  `json:"code"`    // Business response status code
	Message string                  `json:"message"` // Message prompt
	Data    []*models.ApiPostDetail `json:"data"`    // Data
}
