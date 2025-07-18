package common

import "time"

type ContextKey string

const (
	// Context
	REQUEST_TIMESTAMP_MS_CONTEXT_KEY ContextKey = "requestTimestampMS"
	SESSION_TOKEN_CONTEXT_KEY        ContextKey = "sessionToken"
	USER_ID_CONTEXT_KEY              ContextKey = "userID"

	// Environment
	ENVIRONMENT_KEY  string = "ENV"
	DEV_ENVIRONMENT  string = "development"
	PROD_ENVIRONMENT string = "production"

	// Server
	RPC_PORT_KEY          string = "RPC_PORT"
	HTTP_PORT_KEY         string = "HTTP_PORT"
	IDLE_TIMEOUT_SEC_KEY  string = "IDLE_TIMEOUT"
	READ_TIMEOUT_SEC_KEY  string = "READ_TIMEOUT"
	WRITE_TIMEOUT_SEC_KEY string = "WRITE_TIMEOUT"

	// Message Queue
	MESSAGE_QUEUE_HOST_KEY                string = "MESSAGE_QUEUE_HOST"
	MESSAGE_QUEUE_RECONNECT_DELAY_SEC_KEY string = "MESSAGE_QUEUE_RECONNECT_DELAY_SEC"
	MESSAGE_QUEUE_REINIT_DELAY_SEC_KEY    string = "MESSAGE_QUEUE_REINIT_DELAY_SEC"
	MESSAGE_QUEUE_RESEND_DELAY_SEC_KEY    string = "MESSAGE_QUEUE_RESEND_DELAY_SEC"

	// Queue Names
	REVIEW_QUEUE string = "review"

	// Database
	DB_DSN_KEY               string = "DB_DSN"
	DB_MAX_OPEN_CONN_KEY     string = "DB_MAX_OPEN_CONN"
	DB_MAX_IDLE_CONN_KEY     string = "DB_MAX_IDLE_CONN"
	DB_MAX_IDLE_TIME_SEC_KEY string = "DB_MAX_IDLE_TIME_SEC"

	// Object Storage
	OBJECT_STORAGE_ACCESS_KEY     string        = "OBJECT_STORAGE_ACCESS_KEY"
	OBJECT_STORAGE_SECRET_KEY     string        = "OBJECT_STORAGE_SECRET_KEY"
	OBJECT_STORAGE_ENDPOINT_KEY   string        = "OBJECT_STORAGE_ENDPOINT"
	OBJECT_STORAGE_BUCKET_KEY     string        = "OBJECT_STORAGE_BUCKET"
	OBJECT_STORAGE_REGION_KEY     string        = "OBJECT_STORAGE_REGION"
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

	// Constants
	AUTHORIZATION string = "Authorization"
	CONTENT_TYPE  string = "Content-Type"
	ACCEPT        string = "Accept"
)
