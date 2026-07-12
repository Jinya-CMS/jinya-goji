package storage

import (
	"bytes"
	"os"
	"path/filepath"

	"jinya.de/goji/database"
)

func SaveProfilePicture(artistId int, data []byte) error {
	artist, err := database.FindArtistById(artistId)
	if err != nil {
		return err
	}

	err = EnsureSavePathExists()
	if err != nil {
		return err
	}

	hash, err := GetFileHash(bytes.NewReader(data))
	if err != nil {
		return err
	}

	fileName := hash
	finalPath := filepath.Join(SavePath, fileName)

	err = os.WriteFile(finalPath, data, 0644)
	if err != nil {
		return err
	}

	artist.ProfilePicture = new(WebPath + fileName)
	return artist.Update()
}

func DeleteProfilePicture(artistId int) error {
	artist, err := database.FindArtistById(artistId)
	if err != nil {
		return err
	}

	if artist.ProfilePicture != nil && *artist.ProfilePicture != "" {
		// PublicPath is ./public/, artist.ProfilePicture is /jinya-content/hash
		// filepath.Join will handle the / at the start of ProfilePicture correctly if we are careful.
		// In Go, Join("a", "/b") results in "a/b" on Unix.
		fullPath := filepath.Join(PublicPath, *artist.ProfilePicture)
		_ = os.Remove(fullPath)

		artist.ProfilePicture = nil
		return artist.Update()
	}

	return nil
}
