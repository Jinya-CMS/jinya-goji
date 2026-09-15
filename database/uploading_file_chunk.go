package database

type UploadingFileChunk struct {
	Id              int    `db:"id,primarykey,autoincrement"`
	UploadingFileId string `db:"uploading_file_id"`
	ChunkPath       string `db:"chunk_path"`
	ChunkPosition   int    `db:"chunk_position"`
}

func (chunk *UploadingFileChunk) Create() error {
	return dbMap.Insert(chunk)
}

func (chunk *UploadingFileChunk) Update() error {
	_, err := dbMap.Update(chunk)
	return err
}

func (chunk *UploadingFileChunk) Delete() error {
	_, err := dbMap.Delete(chunk)
	return err
}

func FindUploadingFileChunksByFile(uploadingFileId string) ([]UploadingFileChunk, error) {
	return dbMap.SelectType[UploadingFileChunk](`
select id,
       uploading_file_id,
       chunk_path,
       chunk_position
from uploading_file_chunk
where uploading_file_id = $1
order by chunk_position
`, uploadingFileId)
}
