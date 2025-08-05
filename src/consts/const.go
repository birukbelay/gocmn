package consts

import "fmt"

type EnvVar string

const (
	Environment  = EnvVar("ENVIRONMENT")
	RefreshToken = "refresh-token"
	AccessToken  = "access-token"
)
const ApiV1 = "/api/v1"



type OperationId string

func (o OperationId) Str() string {
	return string(o)
}

// FileSize is put in bytes, so 1024 = 1kb, 1024 *1024 ==1mb
type FileSize int64

func (f FileSize) Mb() string {
	return fmt.Sprint(int(f) / 1024)
}
func (f FileSize) Int() int {
	return int(f)
}
func (f FileSize) Int64() int64 {
	return int64(f)
}

// TODO: replace these with env variables
const (
	PanicFile = "panic.log"

	ImageUploadPath = "public/assets/images"
	FileUploadPath  = "public/assets"

	FileServingUrl = "http://localhost:8001/assets/"
)
const (
	OneMb               = FileSize(1024 * 1024) // 1 mb
	MaxUploadSize       = 30 * OneMb            // 30 mb
	MaxSingleUploadSize = 10 * OneMb            // 10 mb
	MaxImagesNumber     = 2
	DefaultPage         = 0
)
