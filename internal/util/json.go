package util

type JSONPayload map[string]any

func NewJSONPayload() JSONPayload {
	return make(map[string]any)
}

func (j JSONPayload) Add(key string, value any) {
	j[key] = value
}
