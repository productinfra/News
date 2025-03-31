package mysql

import (
	"backend/models"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// GetPostTotalCount Query the total number of posts in the database
func GetPostTotalCount() (count int64, err error) {
	sqlStr := `select count(post_id) from post`
	err = db.Get(&count, sqlStr)
	if err != nil {
		zap.L().Error("db.Get(&count, sqlStr) failed", zap.Error(err))
		return 0, err
	}
	return
}

// GetCommunityPostTotalCount Query the total number of posts in the database by community ID
func GetCommunityPostTotalCount(communityID uint64) (count int64, err error) {
	sqlStr := `select count(post_id) from post where community_id = ?`
	err = db.Get(&count, sqlStr, communityID)
	if err != nil {
		zap.L().Error("db.Get(&count, sqlStr) failed", zap.Error(err))
		return 0, err
	}
	return
}

// CreatePost Create a new post
func CreatePost(post *models.Post) (err error) {

	sqlStr := `insert into post(
	post_id, title, content, author_id, community_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, post.PostID, post.Title,
		post.Content, post.AuthorId, post.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return nil
}

// GetPostByID Query post details by ID
func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, status, create_time, update_time
	from post
	where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(ErrorInvalidID)
		}
		zap.L().Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		return nil, errors.New(ErrorQueryFailed)
	}
	return
}

// GetPostListByIDs Query posts by a given list of IDs
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)`
	// Dynamically populate the ids
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	// sqlx.In returns a query with `?` bindvar, we use Rebind() to rebind it
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

// GetPostList Get the list of posts
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	ORDER BY create_time
	DESC 
	limit ?,?
	`
	posts = make([]*models.Post, 0, 2) // 0: length  2: capacity
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByKeywords Query the post list by keyword
func GetPostListByKeywords(p *models.ParamPostList) (posts []*models.Post, err error) {
	// Fuzzy search posts by title or content
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where title like ?
	or content like ?
	ORDER BY create_time
	DESC
	limit ?,?
	`
	// %keyword%
	p.Search = "%" + p.Search + "%"
	posts = make([]*models.Post, 0, 2) // 0: length  2: capacity
	err = db.Select(&posts, sqlStr, p.Search, p.Search, (p.Page-1)*p.Size, p.Size)
	return
}

// GetPostListTotalCount Query the total number of posts by keyword
func GetPostListTotalCount(p *models.ParamPostList) (count int64, err error) {
	// Fuzzy search total count by post title or content
	sqlStr := `select count(post_id)
	from post
	where title like ?
	or content like ?
	`
	// %keyword%
	p.Search = "%" + p.Search + "%"
	err = db.Get(&count, sqlStr, p.Search, p.Search)
	return
}
