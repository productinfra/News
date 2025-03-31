package logic

import (
	"backend/dao/mysql"
	"backend/models"
)

// GetCommunityList Queries the list of communities by category
func GetCommunityList() ([]*models.Community, error) {
	// Query the database to find all communities and return them
	return mysql.GetCommunityList()
}

// GetCommunityDetailByID Queries the details of a community by ID
func GetCommunityDetailByID(id uint64) (*models.CommunityDetailRes, error) {
	return mysql.GetCommunityByID(id)
}
