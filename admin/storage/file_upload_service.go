package storage

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"jinya.de/goji/database"
)

func SaveFileChunk(fileId int, position int, data []byte) (*database.UploadingFileChunk, error) {
	uuidBytes := make([]byte, 16)
	_, _ = rand.Read(uuidBytes)
	uuidStr := fmt.Sprintf("%x-%x-%x-%x-%x", uuidBytes[0:4], uuidBytes[4:6], uuidBytes[6:8], uuidBytes[8:10], uuidBytes[10:])

	tempPath := filepath.Join(os.TempDir(), uuidStr)
	err := os.WriteFile(tempPath, data, 0644)
	if err != nil {
		return nil, err
	}

	uploadingFile, err := database.FindUploadingFileByFile(fileId)
	if err != nil {
		return nil, err
	}

	chunk := database.UploadingFileChunk{
		ChunkPath:       tempPath,
		ChunkPosition:   position,
		UploadingFileId: uploadingFile.Id,
	}

	err = database.Insert(&chunk)
	if err != nil {
		return nil, err
	}

	return &chunk, nil
}

func FinishFileUpload(fileId int) (*database.File, error) {
	file, err := database.FindFileById(fileId)
	if err != nil {
		return nil, err
	}

	uploadingFile, err := database.FindUploadingFileByFile(fileId)
	if err != nil {
		return nil, err
	}

	chunks, err := database.FindUploadingFileChunksByFile(uploadingFile.Id)
	if err != nil {
		return nil, err
	}

	err = EnsureSavePathExists()
	if err != nil {
		return nil, err
	}

	outTempFile, err := os.CreateTemp("", "jinya-upload-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(outTempFile.Name())
	defer outTempFile.Close()

	for _, chunk := range chunks {
		chunkData, err := os.ReadFile(chunk.ChunkPath)
		if err != nil {
			return nil, err
		}
		_, err = outTempFile.Write(chunkData)
		if err != nil {
			return nil, err
		}
	}

	_, err = outTempFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	hash, err := GetFileHash(outTempFile)
	if err != nil {
		return nil, err
	}

	fileName := hash
	finalPath := filepath.Join(SavePath, fileName)

	// Copy temp file to final destination
	_, err = outTempFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	finalFile, err := os.Create(finalPath)
	if err != nil {
		return nil, err
	}
	defer finalFile.Close()
	_, err = io.Copy(finalFile, outTempFile)
	if err != nil {
		return nil, err
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
	file.Path = WebPath + fileName
	err = file.Update()
	if err != nil {
		return nil, err
	}

	err = ClearFileChunks(fileId)
	if err != nil {
		log.Printf("Failed to clear chunks: %v", err)
	}

	err = uploadingFile.Delete()
	if err != nil {
		log.Printf("Failed to delete uploading file record: %v", err)
	}

	return file, nil
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
