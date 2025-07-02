package util

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
