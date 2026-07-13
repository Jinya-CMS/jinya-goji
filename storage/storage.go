package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func getFileHash(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ensureSavePathExists() error {
	return os.MkdirAll(GetStoragePath(), 0775)
}

func GetStoragePath() string {
	path := os.Getenv("STORAGE_PATH")
	if path == "" {
		path = "jinya-content"
	}

	return path
}
