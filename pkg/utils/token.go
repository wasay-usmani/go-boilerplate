package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

func NewUniqueToken(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, GenerateUUID())
}

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateUUIDWithoutHyphen() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	src = rand.NewSource(time.Now().UnixNano())
	r   = rand.New(src) // #nosec G404
)

func StringToken(n int) string {
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[r.Intn(len(letters))]
	}

	return string(result)
}
