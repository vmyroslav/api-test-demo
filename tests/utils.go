package tests

import (
	"math/rand"
	"time"
)

// ToPtr converts a value of any type to a pointer.
func ToPtr[T any](v T) *T {
	return &v
}

// RandomString generates a random string of the given length.
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}
