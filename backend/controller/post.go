package controller

import (
	"backend/logic"
	"backend/models"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler creates a post
func CreatePostHandler(c *gin.Context) {
	// 1. Retrieve and validate parameters
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil { // validator --> binding tag
		zap.L().Debug("c.ShouldBindJSON(post) err", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}

	// Get author ID, which is the current request's UserID (retrieved from c)
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	post.AuthorId = userID

	// 2. Create post
	err = logic.CreatePost(&post)
	if err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. Return response
	ResponseSuccess(c, nil)
}

// PostListHandler retrieves a paginated list of posts
func PostListHandler(c *gin.Context) {
	// Get pagination parameters
	page, size := getPageInfo(c)

	// Retrieve data
	data, err := logic.GetPostList(page, size)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// PostList2Handler retrieves an upgraded post list sorted by creation time or score
func PostList2Handler(c *gin.Context) {
	// GET request parameters (query string): /api/v1/posts2?page=1&size=10&order=time
	// Get pagination parameters
	p := &models.ParamPostList{}

	// c.ShouldBind() selects the appropriate method to retrieve data based on the request data type
	// c.ShouldBindJSON() should be used only when the request contains JSON data
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("PostList2Handler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// Retrieve data
	data, err := logic.GetPostListNew(p) // Updated: unified approach
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// PostDetailHandler retrieves post details by ID
func PostDetailHandler(c *gin.Context) {
	// 1. Retrieve parameters (get post ID from URL)
	postIdStr := c.Param("id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
	}

	// 2. Retrieve post data from database by ID
	post, err := logic.GetPostById(postId)
	if err != nil {
		zap.L().Error("logic.GetPost(postID) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}

	// 3. Return response
	ResponseSuccess(c, post)
}

// GetCommunityPostListHandler retrieves a list of posts based on the community
func GetCommunityPostListHandler(c *gin.Context) {
	// GET request parameters (query string): /api/v1/posts2?page=1&size=10&order=time
	// Get pagination parameters
	p := &models.ParamPostList{
		CommunityID: 0,
		Page:        1,
		Size:        10,
		Order:       models.OrderScore,
	}

	// c.ShouldBind() selects the appropriate method to retrieve data based on the request data type
	// c.ShouldBindJSON() should be used only when the request contains JSON data
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// Retrieve data
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// PostSearchHandler handles post search functionality
func PostSearchHandler(c *gin.Context) {
	p := &models.ParamPostList{}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("PostSearchHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	fmt.Println("Search", p.Search)
	fmt.Println("Order", p.Order)

	// Retrieve data
	data, err := logic.PostSearch(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
