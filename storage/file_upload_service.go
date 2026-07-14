package storage

import (
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"jinya.de/goji/database"
)

func StartFileUpload(fileId int) error {
	uploadingFile := database.UploadingFile{
		Id:     uuid.NewString(),
		FileId: fileId,
	}
	return uploadingFile.Create()
}

func SaveFileChunk(fileId int, position int, data []byte) error {
	tempPath := filepath.Join(os.TempDir(), uuid.New().String())
	err := os.WriteFile(tempPath, data, 0644)
	if err != nil {
		return err
	}

	uploadingFile, err := database.FindUploadingFileByFile(fileId)
	if err != nil {
		return err
	}

	chunk := database.UploadingFileChunk{
		ChunkPath:       tempPath,
		ChunkPosition:   position,
		UploadingFileId: uploadingFile.Id,
	}

	return chunk.Create()
}

func FinishFileUpload(fileId int) error {
	file, err := database.FindFileById(fileId)
	if err != nil {
		return err
	}

	uploadingFile, err := database.FindUploadingFileByFile(fileId)
	if err != nil {
		return err
	}

	chunks, err := database.FindUploadingFileChunksByFile(uploadingFile.Id)
	if err != nil {
		return err
	}

	err = ensureSavePathExists()
	if err != nil {
		return err
	}

	outTempFile, err := os.CreateTemp("", "jinya-upload-*")
	if err != nil {
		return err
	}

	defer os.Remove(outTempFile.Name())
	defer outTempFile.Close()

	for _, chunk := range chunks {
		chunkData, err := os.ReadFile(chunk.ChunkPath)
		if err != nil {
			return err
		}

		_, err = outTempFile.Write(chunkData)
		if err != nil {
			return err
		}
	}

	_, err = outTempFile.Seek(0, 0)
	if err != nil {
		return err
	}

	hash, err := getFileHash(outTempFile)
	if err != nil {
		return err
	}

	fileName := hash
	finalPath := filepath.Join(GetStoragePath(), fileName)

	// Copy temp file to final destination
	_, err = outTempFile.Seek(0, 0)
	if err != nil {
		return err
	}
	finalFile, err := os.Create(finalPath)

	if err != nil {
		return err
	}
	defer finalFile.Close()

	_, err = io.Copy(finalFile, outTempFile)
	if err != nil {
		return err
	}

	// Detect mime type
	// Read first 512 bytes for DetectContentType
	buffer := make([]byte, 512)
	_, _ = outTempFile.Seek(0, 0)
	n, _ := outTempFile.Read(buffer)
	contentType := http.DetectContentType(buffer[:n])
	if contentType == "application/octet-stream" {
		// Fallback to extension based detection if possible
		ext := filepath.Ext(file.Name)
		if ext != "" {
			contentType = mime.TypeByExtension(ext)
			if contentType == "" {
				contentType = "application/octet-stream"
			}
		}
	}

	file.Type = contentType
	file.Path = "/jinya-content/" + fileName
	err = file.Update()
	if err != nil {
		return err
	}

	err = ClearFileChunks(fileId)
	if err != nil {
		slog.Error("Failed to clear chunks", "error", err)
	}

	err = uploadingFile.Delete()
	if err != nil {
		slog.Error("Failed to delete uploading file record", "error", err)
	}

	return nil
}

func ClearFileChunks(fileId int) error {
	uploadingFile, err := database.FindUploadingFileByFile(fileId)
	if err != nil {
		return err
	}

	chunks, err := database.FindUploadingFileChunksByFile(uploadingFile.Id)
	if err != nil {
		return err
	}

	for _, chunk := range chunks {
		_ = os.Remove(chunk.ChunkPath)
		_ = chunk.Delete()
	}

	return nil
}

func DeleteFile(filePath string) error {
	_, file := path.Split(filePath)

	return os.Remove(path.Join(GetStoragePath(), file))
}
