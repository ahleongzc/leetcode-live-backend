package common

import "time"

const (
	ENVIRONMENT_KEY string = "ENV"

	// Server
	PORT_KEY              string = "PORT"
	IDLE_TIMEOUT_SEC_KEY  string = "IDLE_TIMEOUT"
	READ_TIMEOUT_SEC_KEY  string = "READ_TIMEOUT"
	WRITE_TIMEOUT_SEC_KEY string = "WRITE_TIMEOUT"

	// Database
	DB_DSN_KEY               string = "POSTGRES_DSN"
	DB_MAX_OPEN_CONN_KEY     string = "POSTGRES_MAX_OPEN_CONN"
	DB_MAX_IDLE_CONN_KEY     string = "POSTGRES_MAX_IDLE_CONN"
	DB_MAX_IDLE_TIME_SEC_KEY string = "POSTGRES_MAX_IDLE_TIME_SEC"

	// Cloudflare R2
	R2_ACCESS_KEY   string = "R2_ACCESS_KEY"
	R2_SECRET_KEY   string = "R2_SECRET_KEY"
	R2_ENDPOINT_KEY string = "R2_ENDPOINT"
	R2_BUCKET_KEY   string = "R2_BUCKET"
	R2_REGION_KEY   string = "R2_REGION"

	DEV_ENVIRONMENT  string = "development"
	PROD_ENVIRONMENT string = "production"

	DB_QUERY_TIMEOUT time.Duration = 1 * time.Second
)
