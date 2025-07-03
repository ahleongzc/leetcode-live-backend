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
	HOUSEKEEPING_INTERVAL time.Duration = 5 * time.Second
	WRITE_TO_FILE_TIMEOUT time.Duration = 5 * time.Second
	IN_MEMORY_QUEUE_SIZE  uint          = 100
	WORKER_POOL_SIZE      uint          = 20

	// Database
	DB_DSN_KEY               string        = "DB_DSN"
	DB_MAX_OPEN_CONN_KEY     string        = "DB_MAX_OPEN_CONN"
	DB_MAX_IDLE_CONN_KEY     string        = "DB_MAX_IDLE_CONN"
	DB_MAX_IDLE_TIME_SEC_KEY string        = "DB_MAX_IDLE_TIME_SEC"
	DB_QUERY_TIMEOUT         time.Duration = 1 * time.Second

	// Object Storage
	OBJECT_STORAGE_ACCESS_KEY     string        = "OBJECT_STORAGE_ACCESS_KEY"
	OBJECT_STORAGE_SECRET_KEY     string        = "OBJECT_STORAGE_SECRET_KEY"
	OBJECT_STORAGE_ENDPOINT_KEY   string        = "OBJECT_STORAGE_ENDPOINT"
	OBJECT_STORAGE_BUCKET_KEY     string        = "OBJECT_STORAGE_BUCKET"
	OBJECT_STORAGE_REGION_KEY     string        = "OBJECT_STORAGE_REGION"
	FILE_UPLOAD_TIMEOUT           time.Duration = 10 * time.Second
	PRESIGNED_URL_EXPIRY_DURATION time.Duration = 15 * time.Minute

	// AI Providers
	OLLAMA string = "ollama"
	OPENAI string = "openai"

	// LLM
	LLM_PROVIDER_KEY string = "LLM_PROVIDER"
	LLM_MODEL_KEY    string = "LLM_MODEL"
	LLM_BASE_URL_KEY string = "LLM_BASE_URL"
	LLM_API_KEY      string = "LLM_API_KEY"

	// TTS
	TTS_PROVIDER_KEY string = "TTS_PROVIDER"
	TTS_MODEL_KEY    string = "TTS_MODEL"
	TTS_VOICE_KEY    string = "TTS_VOICE"
	TTS_BASE_URL_KEY string = "TTS_BASE_URL"
	TTS_API_KEY      string = "TTS_API_KEY"
	TTS_LANGUAGE_KEY string = "TTS_LANGUAGE"

	// Tables
	SESSION_TABLE_NAME    string = "sessions"
	USER_TABLE_NAME       string = "users"
	TRANSCRIPT_TABLE_NAME string = "transcripts"
	INTERVIEW_TABLE_NAME  string = "interviews"
	QUESTION_TABLE_NAME   string = "questions"

	// HTTP
	HTTP_REQUEST_TIMEOUT       time.Duration = time.Minute
	SESSION_ID_HEADER_KEY      string        = "X-Session-ID"
	AUTHORIZATION              string        = "Authorization"
	CONTENT_TYPE               string        = "Content-Type"
	ACCEPT                     string        = "Accept"
	INCOMING_PAYLOAD_MAX_BYTES int           = 1_048_576
)

var (
	PROD_TRUSTED_ORIGINS = map[string]struct{}{
		"https://leetcode.com": {},
	}
)
