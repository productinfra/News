package controller

import (
	"backend/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Community

// CommunityHandler retrieves the community list
func CommunityHandler(c *gin.Context) {
	// Query all communities (community_id, community_name) and return them as a list
	communityList, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // Do not easily expose server-side errors
		return
	}
	ResponseSuccess(c, communityList)
}

// CommunityDetailHandler retrieves community details by ID
func CommunityDetailHandler(c *gin.Context) {
	// 1. Get community ID
	communityIdStr := c.Param("id")                               // Get URL parameter
	communityId, err := strconv.ParseUint(communityIdStr, 10, 64) // Convert ID string format
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2. Retrieve community details by ID
	communityList, err := logic.GetCommunityDetailByID(communityId)
	if err != nil {
		zap.L().Error("logic.GetCommunityByID() failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeSuccess, err.Error())
		return
	}
	ResponseSuccess(c, communityList)
}
