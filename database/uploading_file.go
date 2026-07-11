package database

type UploadingFile struct {
	Id     string `db:"id,primarykey" json:"id"`
	FileId int    `db:"file_id" json:"fileId"`
}
