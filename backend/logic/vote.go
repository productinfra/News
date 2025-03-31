package logic

import (
	"backend/dao/redis"
	"backend/models"
	"strconv"

	"go.uber.org/zap"
)

// 1. User vote data

/*
Voting algorithm: http://www.ruanyifeng.com/blog/2012/03/ranking_algorithm_reddit.html
This project uses a simplified version of the voting score.
Each vote adds 432 points, 86400/200 -> 200 upvotes will allow the post to stay on the homepage for another day -> "Redis in Action"
*/

/* PostVote - Voting for a post
There are four voting scenarios:
1. Upvote (1)
2. Downvote (-1)
3. Cancel vote (0)
4. Reverse vote

Record the users who voted on the post.
Update the post score: Upvotes add points; Downvotes subtract points.

When v=1, there are two scenarios:
	1. The user has not voted before and now votes for the post --> Update score and voting record. Absolute value of the difference: 1 +432.
	2. The user previously voted down and now changes to upvote --> Update score and voting record. Absolute value of the difference: 2 +432*2.
When v=0, there are two scenarios:
	1. The user previously voted down and now cancels the vote --> Update score and voting record. Absolute value of the difference: 1 +432.
	2. The user previously voted up and now cancels the vote --> Update score and voting record. Absolute value of the difference: 1 -432.
When v=-1, there are two scenarios:
	1. The user has not voted before and now votes down --> Update score and voting record. Absolute value of the difference: 1 -432.
	2. The user previously voted up and now changes to downvote --> Update score and voting record. Absolute value of the difference: 2 -432*2.

Voting restrictions:
Users can vote on a post within one week of its creation. After one week, voting is not allowed.
	1. After the deadline, the upvote and downvote counts stored in Redis will be saved to the MySQL table.
	2. After the deadline, the KeyPostVotedZSetPrefix will be deleted.
*/

// VoteForPost - Voting functionality to vote for a post
func VoteForPost(userId uint64, p *models.VoteDataForm) error {
	zap.L().Debug("VoteForPost",
		zap.Uint64("userId", userId),
		zap.String("postId", p.PostID),
		zap.Int8("Direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userId)), p.PostID, float64(p.Direction))
}
