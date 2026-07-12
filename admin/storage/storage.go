package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

const (
	WebPath    = "/jinya-content/"
	PublicPath = "./public/"
	SavePath   = PublicPath + WebPath
)

func GetFileHash(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func EnsureSavePathExists() error {
	return os.MkdirAll(SavePath, 0775)
}
