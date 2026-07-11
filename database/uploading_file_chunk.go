package database

type UploadingFileChunk struct {
	Id              int    `db:"id,primarykey,autoincrement" json:"id"`
	UploadingFileId string `db:"uploading_file_id" json:"uploadingFileId"`
	ChunkPath       string `db:"chunk_path" json:"chunkPath"`
	ChunkPosition   int    `db:"chunk_position" json:"chunkPosition"`
}
