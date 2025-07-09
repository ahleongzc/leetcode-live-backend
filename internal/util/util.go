package util

import "strings"

func ToPtr[T any](v T) *T {
	return &v
}

func FromPtr[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

func MillisToSeconds(milliSec int64) int64 {
	return milliSec / 1000
}

func ContainsNewline(text string) bool {
	return strings.Contains(text, "\n")
}
