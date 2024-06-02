package internal

import (
	"os"
)

func getenv[T any](key string, fallback T) any {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
