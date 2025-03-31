package logic

import (
	"backend/dao/mysql"
	"backend/dao/redis"
	"backend/models"
	"backend/pkg/snowflake"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

// CreatePost Creates a post
func CreatePost(post *models.Post) (err error) {
	// 1. Generate post_id (create post ID)
	postID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		return
	}
	post.PostID = postID
	// 2. Create the post and save it to the database
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		return err
	}
	community, err := mysql.GetCommunityNameByID(fmt.Sprint(post.CommunityID))
	if err != nil {
		zap.L().Error("mysql.GetCommunityNameByID failed", zap.Error(err))
		return err
	}
	// Store post information in Redis
	if err := redis.CreatePost(
		post.PostID,
		post.AuthorId,
		post.Title,
		TruncateByWords(post.Content, 120),
		community.CommunityID); err != nil {
		zap.L().Error("redis.CreatePost failed", zap.Error(err))
		return err
	}
	return
}

// GetPostById Fetch post details by ID
func GetPostById(postID int64) (data *models.ApiPostDetail, err error) {
	// Query and assemble the data we want to use for the interface
	// Query post information
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(postID) failed",
			zap.Int64("postID", postID),
			zap.Error(err))
		return nil, err
	}
	// Query author information by author id
	user, err := mysql.GetUserByID(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed",
			zap.Uint64("postID", post.AuthorId),
			zap.Error(err))
		return
	}
	// Query community details by community id
	community, err := mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed",
			zap.Uint64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	// Query vote count for the post by post id
	voteNum, err := redis.GetPostVoteNum(postID)

	// Assemble data for the interface
	data = &models.ApiPostDetail{
		Post:               post,
		CommunityDetailRes: community,
		AuthorName:         user.UserName,
		VoteNum:            voteNum,
	}
	return data, nil
}

// GetPostList Fetch post list
func GetPostList(page, size int64) ([]*models.ApiPostDetail, error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed")
		return nil, err
	}
	data := make([]*models.ApiPostDetail, 0, len(postList)) // Initialize data
	for _, post := range postList {
		// Query author information by author id
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", post.AuthorId),
				zap.Error(err))
			continue
		}
		// Query community details by community id
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		// Assemble data for the interface
		postDetail := &models.ApiPostDetail{
			Post:               post,
			CommunityDetailRes: community,
			AuthorName:         user.UserName,
		}
		data = append(data, postDetail)
	}
	return data, nil
}

// GetPostList2 Upgraded version of the post list interface: sorted by creation time or score
func GetPostList2(p *models.ParamPostList) (*models.ApiPostDetailRes, error) {
	var res models.ApiPostDetailRes
	// Get total number of posts from MySQL
	total, err := mysql.GetPostTotalCount()
	if err != nil {
		return nil, err
	}
	res.Page.Total = total
	// 1. Query post IDs from Redis according to the sorting rule in the parameters
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) returned 0 data")
		return &res, nil
	}
	zap.L().Debug("GetPostList2", zap.Any("ids: ", ids))
	// 2. Pre-fetch vote count data for each post
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	// 3. Query post details from the database using the post IDs
	// Ensure the data is returned in the same order as the post IDs: order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}
	res.Page.Page = p.Page
	res.Page.Size = p.Size
	res.List = make([]*models.ApiPostDetail, 0, len(posts))
	// 4. Assemble the data
	// Query the author and community information for each post and fill it in
	for idx, post := range posts {
		// Query author information by author id
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", post.AuthorId),
				zap.Error(err))
			user = nil
		}
		// Query community details by community id
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", post.CommunityID),
				zap.Error(err))
			community = nil
		}
		// Assemble data for the interface
		postDetail := &models.ApiPostDetail{
			VoteNum:            voteData[idx],
			Post:               post,
			CommunityDetailRes: community,
			AuthorName:         user.UserName,
		}
		res.List = append(res.List, postDetail)
	}
	return &res, nil
}

// GetCommunityPostList Query post list by community id
func GetCommunityPostList(p *models.ParamPostList) (*models.ApiPostDetailRes, error) {
	var res models.ApiPostDetailRes
	// Get total number of posts for the community from MySQL
	total, err := mysql.GetCommunityPostTotalCount(p.CommunityID)
	if err != nil {
		return nil, err
	}
	res.Page.Total = total
	// 1. Query post IDs from Redis based on the sorting rule in the parameters
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostList(p) returned 0 data")
		return &res, nil
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 2. Pre-fetch vote count data for each post
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	// 3. Query post details from the database using the post IDs
	// Ensure the data is returned in the same order as the post IDs: order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}
	res.Page.Page = p.Page
	res.Page.Size = p.Size
	res.List = make([]*models.ApiPostDetail, 0, len(posts))
	// 4. Query community details by community id
	// To reduce database queries, we fetch community info in advance
	community, err := mysql.GetCommunityByID(p.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed",
			zap.Uint64("community_id", p.CommunityID),
			zap.Error(err))
		community = nil
	}
	for idx, post := range posts {
		// Filter out posts that do not belong to the current community
		if post.CommunityID != p.CommunityID {
			continue
		}
		// Query author information by author id
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", post.AuthorId),
				zap.Error(err))
			user = nil
		}
		// Assemble data for the interface
		postDetail := &models.ApiPostDetail{
			VoteNum:            voteData[idx],
			Post:               post,
			CommunityDetailRes: community,
			AuthorName:         user.UserName,
		}
		res.List = append(res.List, postDetail)
	}
	return &res, nil
}

// GetPostListNew Combines two post list query logics into one function
func GetPostListNew(p *models.ParamPostList) (data *models.ApiPostDetailRes, err error) {
	// Execute different business logic based on the request parameters
	if p.CommunityID == 0 {
		// Query all posts
		data, err = GetPostList2(p)
	} else {
		// Query posts by community id
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return data, nil
}

// PostSearch Search posts by keywords
func PostSearch(p *models.ParamPostList) (*models.ApiPostDetailRes, error) {
	var res models.ApiPostDetailRes
	// Get total number of posts from MySQL based on search criteria
	total, err := mysql.GetPostListTotalCount(p)
	if err != nil {
		return nil, err
	}
	res.Page.Total = total
	// 1. Query posts matching the search criteria from MySQL
	posts, err := mysql.GetPostListByKeywords(p)
	if err != nil {
		return nil, err
	}
	// If no posts are found
	if len(posts) == 0 {
		return &models.ApiPostDetailRes{}, nil
	}
	// 2. Query the vote count for the posts from Redis
	ids := make([]string, 0, len(posts))
	for _, post := range posts {
		ids = append(ids, strconv.Itoa(int(post.PostID)))
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	res.Page.Size = p.Size
	res.Page.Page = p.Page
	// 3. Assemble data
	res.List = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		// Query author information by author id
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID", post.AuthorId),
				zap.Error(err))
			user = nil
		}
		// Query community details by community id
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", post.CommunityID),
				zap.Error(err))
			community = nil
		}
		// Assemble data for the interface
		postDetail := &models.ApiPostDetail{
			VoteNum:            voteData[idx],
			Post:               post,
			CommunityDetailRes: community,
			AuthorName:         user.UserName,
		}
		res.List = append(res.List, postDetail)
	}
	return &res, nil
}
