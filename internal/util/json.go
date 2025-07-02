package util

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ahleongzc/leetcode-live-backend/internal/common"
)

type JSONPayload map[string]any

func NewJSONPayload() JSONPayload {
	return make(map[string]any)
}

func (j JSONPayload) Add(key string, value any) {
	j[key] = value
}

func StringToJSON(input string, dst any) error {
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")

	if start == -1 || end == -1 || start > end {
		return fmt.Errorf("invalid JSON format: missing or misplaced '{' or '}', the string is %s: %w", input, common.ErrInternalServerError)
	}

	jsonString := input[start : end+1]
	if err := json.Unmarshal([]byte(jsonString), dst); err != nil {
		return fmt.Errorf("failed to unmarshal JSON, %s: %w", err, common.ErrInternalServerError)
	}

	return nil
}
