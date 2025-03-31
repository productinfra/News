package redis

import (
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	OneWeekInSeconds          = 7 * 24 * 3600        // One week in seconds
	OneMonthInSeconds         = 4 * OneWeekInSeconds // One month in seconds
	VoteScore         float64 = 432                  // Value of each vote (432 points)
	PostPerAge                = 20                   // Show 20 posts per page
)

/*
Voting Algorithm: http://www.ruanyifeng.com/blog/2012/03/ranking_algorithm_reddit.html
This project uses a simplified version of the voting score.
Each vote adds 432 points to the score. 86400/200 -> 200 upvotes can keep the post on the homepage for another day -> "Redis in Action"
*/

/* PostVote for voting on a post
There are four cases for voting: 1. Upvote 2. Downvote 3. Cancel vote 4. Reverse vote

Record users who participated in voting
Update the post score: upvotes add score; downvotes subtract score

v=1 when there are two situations:
	1. Previously no vote, now upvoting --> Update the score and voting record. The absolute difference: 1 +432
	2. Previously downvoted, now changing to upvote --> Update the score and voting record. The absolute difference: 2 +432*2
v=0 when there are two situations:
	1. Previously downvoted, now canceling vote --> Update the score and voting record. The absolute difference: 1 +432
	2. Previously upvoted, now canceling vote --> Update the score and voting record. The absolute difference: 1 -432
v=-1 when there are two situations:
	1. Previously no vote, now downvoting --> Update the score and voting record. The absolute difference: 1 -432
	2. Previously upvoted, now changing to downvote --> Update the score and voting record. The absolute difference: 2 -432*2

Voting restrictions:
Users can vote within one week of the post being published; after one week, voting is not allowed.
	1. After the expiration, save the upvotes and downvotes from Redis to the MySQL table
	2. After the expiration, delete the KeyPostVotedZSetPrefix
*/

// VoteForPost Votes for a post
func VoteForPost(userID string, postID string, v float64) (err error) {
	// 1. Check voting restrictions
	// Get the post's creation time from Redis
	postTime := client.ZScore(KeyPostTimeZSet, postID).Val()
	if float64(time.Now().Unix())-postTime > OneWeekInSeconds { // Voting is not allowed after one week
		// Voting is not allowed
		return ErrorVoteTimeExpire
	}
	// 2. Update the post's score
	// Steps 2 and 3 need to be put into a pipeline transaction
	// Check if the user has voted by checking the current voting record for the post
	key := KeyPostVotedZSetPrefix + postID
	ov := client.ZScore(key, userID).Val()

	// Update: If the value of the vote this time is the same as the previously saved value, disallow repeat voting
	if v == ov {
		return ErrVoteRepeated
	}
	var op float64
	if v > ov {
		op = 1
	} else {
		op = -1
	}
	diffAbs := math.Abs(ov - v)                // Calculate the absolute difference between the two votes
	pipeline := client.TxPipeline()            // Pipeline transaction
	incrementScore := VoteScore * diffAbs * op // Calculate the score change (addition)
	// ZIncrBy is used to increase the score of a member in a sorted set by a specified amount
	_, err = pipeline.ZIncrBy(KeyPostScoreZSet, incrementScore, postID).Result() // Update the score
	if err != nil {
		return err
	}
	// 3. Record the user's vote for the post
	if v == 0 {
		_, err = client.ZRem(key, postID).Result()
	} else {
		pipeline.ZAdd(key, redis.Z{ // Record the vote
			Score:  v, // Upvote or downvote
			Member: userID,
		})
	}
	// 4. Update the post's vote count
	pipeline.HIncrBy(KeyPostInfoHashPrefix+postID, "votes", int64(op))

	_, err = pipeline.Exec()
	return err
}

// CreatePost Stores post information in Redis using hash
func CreatePost(postID, userID uint64, title, summary string, CommunityID uint64) (err error) {
	now := float64(time.Now().Unix())
	votedKey := KeyPostVotedZSetPrefix + strconv.Itoa(int(postID))
	communityKey := KeyCommunityPostSetPrefix + strconv.Itoa(int(CommunityID))
	postInfo := map[string]interface{}{
		"title":    title,
		"summary":  summary,
		"post:id":  postID,
		"user:id":  userID,
		"time":     now,
		"votes":    1,
		"comments": 0,
	}

	// Pipeline transaction
	pipeline := client.TxPipeline()
	// Voting zSet
	pipeline.ZAdd(votedKey, redis.Z{ // The author automatically votes for the post
		Score:  1,
		Member: userID,
	})
	pipeline.Expire(votedKey, time.Second*OneMonthInSeconds*6) // Expiry time: 6 months
	// Post hash
	pipeline.HMSet(KeyPostInfoHashPrefix+strconv.Itoa(int(postID)), postInfo)
	// Add to the score ZSet
	pipeline.ZAdd(KeyPostScoreZSet, redis.Z{
		Score:  now + VoteScore,
		Member: postID,
	})
	// Add to the time ZSet
	pipeline.ZAdd(KeyPostTimeZSet, redis.Z{
		Score:  now,
		Member: postID,
	})
	// Add to the corresponding community set
	pipeline.SAdd(communityKey, postID)
	_, err = pipeline.Exec()
	return
}

// GetPost Retrieves posts from Redis in a paginated manner
func GetPost(order string, page int64) []map[string]string {
	key := KeyPostScoreZSet
	if order == "time" {
		key = KeyPostTimeZSet
	}
	start := (page - 1) * PostPerAge
	end := start + PostPerAge - 1
	ids := client.ZRevRange(key, start, end).Val()
	postList := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		postData := client.HGetAll(KeyPostInfoHashPrefix + id).Val()
		postData["id"] = id
		postList = append(postList, postData)
	}
	return postList
}

// GetCommunityPost Retrieves posts from a specific community by creation time or score with pagination
func GetCommunityPost(communityName, orderKey string, page int64) []map[string]string {
	key := orderKey + communityName // Create cache key

	if client.Exists(key).Val() < 1 {
		client.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, KeyCommunityPostSetPrefix+communityName, orderKey)
		client.Expire(key, 60*time.Second)
	}
	return GetPost(key, page)
}

// Reddit Hot Rank Algorithm
// From https://github.com/reddit-archive/reddit/blob/master/r2/r2/lib/db/_sorts.pyx
func Hot(ups, downs int, date time.Time) float64 {
	s := float64(ups - downs)
	order := math.Log10(math.Max(math.Abs(s), 1))
	var sign float64
	if s > 0 {
		sign = 1
	} else if s == 0 {
		sign = 0
	} else {
		sign = -1
	}
	seconds := float64(date.Second() - 1577808000)
	return math.Round(sign*order + seconds/43200)
}
