package eightid

import (
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// toBase62 converts an int64 to base62 string
func toBase62(n int64) string {
	const base = 62
	const digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append([]byte{digits[n%base]}, result...)
		n /= base
	}
	return string(result)
}

// GenerateID generates a unique, 8-char, sortable, case-sensitive ID
func New() string {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	base62Time := toBase62(now)

	// Last 2 chars from timestamp for sortable prefix
	timePart := base62Time[len(base62Time)-2:]

	// 6 random characters
	var sb strings.Builder
	for i := 0; i < 6; i++ {
		sb.WriteByte(charset[seededRand.Intn(len(charset))])
	}

	return timePart + sb.String()
}

// IsValid checks if ID is 8 chars and matches allowed characters
func IsValid(id string) bool {
	if len(id) != 8 {
		return false
	}
	for _, ch := range id {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9')) {
			return false
		}
	}
	return true
}
