package storage

import (
	"bytes"
	"os"
	"path"

	"jinya.de/goji/database"
)

func SaveProfilePicture(artistId int, data []byte) error {
	artist, err := database.FindArtistById(artistId)
	if err != nil {
		return err
	}

	err = ensureSavePathExists()
	if err != nil {
		return err
	}

	hash, err := getFileHash(bytes.NewReader(data))
	if err != nil {
		return err
	}

	fileName := hash
	finalPath := path.Join(GetStoragePath(), fileName)

	err = os.WriteFile(finalPath, data, 0644)
	if err != nil {
		return err
	}

	artist.ProfilePicture = new("/jinya-content/" + fileName)
	return artist.Update()
}

func DeleteProfilePicture(artistId int) error {
	artist, err := database.FindArtistById(artistId)
	if err != nil {
		return err
	}

	if artist.ProfilePicture != nil && *artist.ProfilePicture != "" {
		_, file := path.Split(*artist.ProfilePicture)
		fullPath := path.Join(GetStoragePath(), file)
		_ = os.Remove(fullPath)

		artist.ProfilePicture = nil
		return artist.Update()
	}

	return nil
}
