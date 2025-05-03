package helper

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(str string) string {
	hash := sha256.Sum256([]byte(str))

	return hex.EncodeToString(hash[:])
}
