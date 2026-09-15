package database

type UploadingFile struct {
	Id     string `db:"id,primarykey"`
	FileId int    `db:"file_id,unique"`
}

func (file *UploadingFile) Create() error {
	return dbMap.Insert(file)
}

func (file *UploadingFile) Update() error {
	_, err := dbMap.Update(file)
	return err
}

func (file *UploadingFile) Delete() error {
	_, err := dbMap.Delete(file)
	return err
}

func FindUploadingFileByFile(fileId int) (*UploadingFile, error) {
	return dbMap.SelectOneType[UploadingFile](`
select id,
       file_id
from uploading_file
where file_id = $1
`, fileId)
}

func (file *UploadingFile) GetChunks() ([]UploadingFileChunk, error) {
	return FindUploadingFileChunksByFile(file.Id)
}

func (file *UploadingFile) GetFile() (*File, error) {
	return FindFileById(file.FileId)
}
