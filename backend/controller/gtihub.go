package controller

import (
	"backend/logic"
	"backend/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GithubTrendingHandler retrieves trending GitHub projects
func GithubTrendingHandler(c *gin.Context) {
	p := &models.ParamGithubTrending{}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GithubTrendingHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// Fetch data
	data, err := logic.GetGithubTrending(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
