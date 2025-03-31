package redis

import (
	"backend/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// getIDsFromKey Query a specified number of elements in descending order of score
func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZRevRange: Query elements in descending order of score
	return client.ZRevRange(key, start, end).Result()
}

// GetPostIDsInOrder Upgraded voting list interface: Sort by creation time or by score (the queried IDs are already sorted from high to low based on the order)
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// Retrieve IDs from Redis
	// 1. Determine the Redis key to query based on the order parameter in the user's request
	key := KeyPostTimeZSet            // Default is time
	if p.Order == models.OrderScore { // Request by score
		key = KeyPostScoreZSet
	}
	// 2. Determine the starting point for the query index
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData Query the voting data of posts based on IDs (get the number of upvotes for each post)
func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := KeyPostVotedZSetPrefix + id
		// Find the count of elements with a score of 1 -> Count the number of upvotes for each post
		v := client.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}
	// Use pipeline to send multiple commands at once to reduce RTT
	// pipeline := client.Pipeline()
	// for _, id := range ids {
	// 	key := KeyCommunityPostSetPrefix + id
	// 	pipeline.ZCount(key, "1", "1") // ZCount returns the count of elements within the min and max score range
	// }
	// cmders, err := pipeline.Exec()
	// if err != nil {
	// 	return nil, err
	// }
	// data = make([]int64, 0, len(cmders))
	// for _, cmder := range cmders {
	// 	v := cmder.(*redis.IntCmd).Val()
	// 	data = append(data, v)
	// }
	return data, nil
}

// GetPostVoteNum Query the voting data of a post based on ID (get the number of upvotes for the post)
func GetPostVoteNum(ids int64) (data int64, err error) {
	key := KeyPostVotedZSetPrefix + strconv.Itoa(int(ids))
	// Find the count of elements with a score of 1 -> Count the number of upvotes for the post
	data = client.ZCount(key, "1", "1").Val()
	return data, nil
}

/*
*
  - @Author huchao
  - @Description //TODO Query community post IDs (the queried IDs are already sorted from high to low based on the order)
  - @Date 23:06 2022/2/16
  - @Param orderKey: Sort by score or time
    Use zinterstore to combine the community key and order key (community or time)
    *
*/
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 1. Determine the Redis key to query based on the order parameter in the user's request
	orderkey := KeyPostTimeZSet       // Default is time
	if p.Order == models.OrderScore { // Request by score
		orderkey = KeyPostScoreZSet
	}

	// Use zinterstore to combine the community post set and the post score zset to generate a new zset
	// For the new zset, query data according to the previous logic

	// Community key
	cKey := KeyCommunityPostSetPrefix + strconv.Itoa(int(p.CommunityID))

	// Use cache key to reduce the number of zinterstore executions
	key := orderkey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		// Not exists, need to compute
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX", // When aggregating two zsets, take the maximum value
		}, cKey, orderkey) // zinterstore calculation
		pipeline.Expire(key, 60*time.Second) // Set expiration time
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// If exists, directly query IDs based on the key
	return getIDsFromKey(key, p.Page, p.Size)
}
