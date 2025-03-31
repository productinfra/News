package mysql

import (
	"backend/models"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// CreateComment creates a new comment in the database
func CreateComment(comment *models.Comment) (err error) {
	sqlStr := `insert into comment(
	comment_id, content, post_id, author_id, parent_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, comment.CommentID, comment.Content, comment.PostID,
		comment.AuthorID, comment.ParentID)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

// GetCommentListByIDs retrieves a list of comments by their IDs
func GetCommentListByIDs(ids []string) (commentList []*models.Comment, err error) {
	sqlStr := `select comment_id, content, post_id, author_id, parent_id, create_time
	from comment
	where comment_id in (?)`
	// Dynamically fill in the IDs
	query, args, err := sqlx.In(sqlStr, ids)
	if err != nil {
		return
	}
	// sqlx.In returns a query with `?` bind variables, we use Rebind() to rebind it
	query = db.Rebind(query)
	err = db.Select(&commentList, query, args...)
	return
}
