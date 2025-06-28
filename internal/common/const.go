package common

import "time"

type ContextKey string

const (
	// Context
	REQUEST_TIMESTAMP_MS_CONTEXT_KEY ContextKey = "requestTimestampMS"

	// Environment
	ENVIRONMENT_KEY  string = "ENV"
	DEV_ENVIRONMENT  string = "development"
	PROD_ENVIRONMENT string = "production"

	// Server
	PORT_KEY              string        = "PORT"
	IDLE_TIMEOUT_SEC_KEY  string        = "IDLE_TIMEOUT"
	READ_TIMEOUT_SEC_KEY  string        = "READ_TIMEOUT"
	WRITE_TIMEOUT_SEC_KEY string        = "WRITE_TIMEOUT"
	HOUSEKEEPING_INTERVAL time.Duration = 30 * time.Second

	// Database
	DB_DSN_KEY               string        = "POSTGRES_DSN"
	DB_MAX_OPEN_CONN_KEY     string        = "POSTGRES_MAX_OPEN_CONN"
	DB_MAX_IDLE_CONN_KEY     string        = "POSTGRES_MAX_IDLE_CONN"
	DB_MAX_IDLE_TIME_SEC_KEY string        = "POSTGRES_MAX_IDLE_TIME_SEC"
	DB_QUERY_TIMEOUT         time.Duration = time.Second

	// Cloudflare R2
	R2_ACCESS_KEY                 string        = "R2_ACCESS_KEY"
	R2_SECRET_KEY                 string        = "R2_SECRET_KEY"
	R2_ENDPOINT_KEY               string        = "R2_ENDPOINT"
	R2_BUCKET_KEY                 string        = "R2_BUCKET"
	R2_REGION_KEY                 string        = "R2_REGION"
	FILE_UPLOAD_TIMEOUT           time.Duration = time.Minute
	PRESIGNED_URL_EXPIRY_DURATION time.Duration = 15 * time.Minute

	// OpenAI
	OPENAI_API_KEY  string = "OPENAI_API_KEY"
	OPENAI_BASE_URL string = "https://api.openai.com"

	// Ollama
	OLLAMA_BASE_URL string = "http://localhost:11434"

	// TTS
	TTS_REQUEST_TIMEOUT time.Duration = 10 * time.Second

	// Tables
	SESSION_TABLE_NAME    string = "sessions"
	USER_TABLE_NAME       string = "users"
	TRANSCRIPT_TABLE_NAME string = "transcripts"

	// HTTP
	HTTP_REQUEST_TIMEOUT       time.Duration = time.Minute
	AUTHORIZATION              string        = "Authorization"
	CONTENT_TYPE               string        = "Content-Type"
	ACCEPT                     string        = "Accept"
	INCOMING_PAYLOAD_MAX_BYTES int           = 1_048_576
)

var (
	DEV_TRUSTED_ORIGINS = map[string]struct{}{
		"localhost": {},
		"0.0.0.0":   {},
		"":          {},
	}

	PROD_TRUSTED_ORIGINS = map[string]struct{}{
		"https://leetcode.com": {},
	}
)
