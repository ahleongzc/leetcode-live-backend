package config

import "time"

const (
	// Pool size
	IN_MEMORY_QUEUE_SIZE                  uint = 100
	CONSUMER_POOL_SIZE                    uint = 20
	INTENT_CLASSIFICATION_MODEL_POOL_SIZE uint = 5

	// Interval
	HOUSEKEEPING_INTERVAL time.Duration = 5 * time.Second

	// Timeout
	DB_QUERY_TIMEOUT                 time.Duration = 1 * time.Second
	PUBLISHER_TIMEOUT                time.Duration = 5 * time.Second
	FILE_UPLOAD_TIMEOUT              time.Duration = 10 * time.Second
	HTTP_REQUEST_TIMEOUT             time.Duration = time.Minute
	WRITE_TO_FILE_TIMEOUT            time.Duration = 5 * time.Second
	MESSAGE_QUEUE_CONNECTION_TIMEOUT time.Duration = 30 * time.Second

	// HTTP
	INCOMING_PAYLOAD_MAX_BYTES int    = 1_048_576
	SESSION_TOKEN_HEADER_KEY   string = "X-Session-Token"
	INTERVIEW_TOKEN_HEADER_KEY string = "X-Interview-Token"

	// Pagination
	PAGINATION_DEFAULT_OFFSET uint = 0
	PAGINATION_DEFAULT_LIMIT  uint = 10
	PAGINATION_MAX_LIMIT      uint = 20
)

var (
	PROD_TRUSTED_ORIGINS = map[string]struct{}{
		"https://leetcode.com": {},
	}
)
