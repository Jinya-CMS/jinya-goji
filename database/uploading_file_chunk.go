package database

type UploadingFileChunk struct {
	Id              int    `db:"id,primarykey,autoincrement"`
	UploadingFileId string `db:"uploading_file_id"`
	ChunkPath       string `db:"chunk_path"`
	ChunkPosition   int    `db:"chunk_position"`
}

func (chunk *UploadingFileChunk) Create() error {
	return Insert(chunk)
}

func (chunk *UploadingFileChunk) Update() error {
	return Update(chunk)
}

func (chunk *UploadingFileChunk) Delete() error {
	return Delete(chunk)
}

func FindUploadingFileChunksByFile(uploadingFileId string) ([]UploadingFileChunk, error) {
	return Select[UploadingFileChunk](`
select id,
       uploading_file_id,
       chunk_path,
       chunk_position
from uploading_file_chunk
where uploading_file_id = $1
order by chunk_position
`, uploadingFileId)
}
