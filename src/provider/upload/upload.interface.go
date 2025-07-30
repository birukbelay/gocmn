package upload

import (
	"mime/multipart"

	"github.com/birukbelay/gocmn/src/dtos"
)

// FileUploadInterface is only concerned with the file upload, not the authentication or database
type FileUploadInterface interface {
	//UploadSingleFile TODO accept allowed FileTypes & prefix as a param
	UploadSingleFile(file *multipart.FileHeader) (dtos.GResp[UploadDto], error)
	UploadWithGivenName(file *multipart.FileHeader, name string) (dtos.GResp[UploadDto], error)
	DeleteFileWithName(fileName string) error
}

type UploadDto struct {
	Name     string `gorm:"index:,unique;not null" json:"name,omitempty"`
	Url      string `json:"url"`
	Path     string `json:"-"`
	UploadId string `json:"uploadId,omitempty"`
	Hash     string `json:"hash"`
	FileType string `json:"file_type"`
	Size     int64  `json:"size"`
	ETag     string
	Ext      string `json:"ext"`
}
