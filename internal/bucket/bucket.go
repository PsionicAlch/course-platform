package bucket

import "embed"

type File struct {
	Name     string
	Checksum string
}

type Bucket interface {
	GetAllFiles() ([]*File, error)
	UploadFileFS(files embed.FS, fileName, checksum string) error
	DeleteFile(fileName string) error
}
