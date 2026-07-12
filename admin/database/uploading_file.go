package database

type UploadingFile struct {
	Id     string `db:"id,primarykey"`
	FileId int    `db:"file_id,unique"`
}

func (file *UploadingFile) Create() error {
	return Insert(file)
}

func (file *UploadingFile) Update() error {
	return Update(file)
}

func (file *UploadingFile) Delete() error {
	return Delete(file)
}

func FindUploadingFileById(id string) (*UploadingFile, error) {
	return Get[UploadingFile](id)
}

func FindUploadingFileByFile(fileId int) (*UploadingFile, error) {
	return SelectOne[UploadingFile](`
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
