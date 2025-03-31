package redis

// Redis keys: Note the use of namespaces for easier querying and splitting
const (
	KeyPostInfoHashPrefix = "news-plus:post:"      // Hash; Post information prefix
	KeyPostTimeZSet       = "news-plus:post:time"  // ZSet; Posts and posting times definition
	KeyPostScoreZSet      = "news-plus:post:score" // ZSet; Posts and voting scores definition
	// KeyPostVotedUpSetPrefix   = "news-plus:post:voted:down:" // Commented out
	// KeyPostVotedDownSetPrefix = "news-plus:post:voted:up:"   // Commented out
	KeyPostVotedZSetPrefix    = "news-plus:post:voted:" // ZSet; Records of users and vote types; parameter is post_id
	KeyCommunityPostSetPrefix = "news-plus:community:"  // Set to store post IDs for each community section
)
