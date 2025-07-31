package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/rs/zerolog/log"

	file_const "github.com/birukbelay/gocmn/src/consts"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/resp_const"
)

// ValidateFile validates & calculate size, detect filetype, calculate hash
func ValidateFile(header *multipart.FileHeader) (dtos.GResp[*UploadDto], error) {

	// =====================       Validate the file Size the hash      =====================
	//==========================================================================
	fileSize := header.Size
	if fileSize > file_const.MaxSingleUploadSize.Int64() {
		return dtos.BadReqRespMsgCode[*UploadDto](resp_const.FileTooBig.Msg(), resp_const.FileTooBig),
			errors.New("single CoverImage size Exeded limit of 10 MB" + fmt.Sprint(fileSize))
	}
	file, err := header.Open()
	defer file.Close()
	if err != nil {
		return dtos.BadReqRespMsgCode[*UploadDto](err.Error(), resp_const.InvalidFileType), err
	}

	val, err := GetFileInfo(file)
	if err == nil {
		val.Body.Size = fileSize
	}
	return val, err
}

// ValidateFile validates & calculate size, detect filetype, calculate hash
func GetFileInfo(file multipart.File) (dtos.GResp[*UploadDto], error) {

	// =====================       Calculate the hash      =====================
	//==========================================================================
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return dtos.BadReqRespMsgCode[*UploadDto](err.Error(), resp_const.InvalidFileType), err
	}
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return dtos.BadReqRespMsgCode[*UploadDto](err.Error(), resp_const.InvalidFileType), err
	}
	/* ==============================  Detect File Type ================================
	====================================================================================*/
	//TODO: make the function take the allowed types and size as a parameter
	contentType, er := DetectFileType(file)
	if er != nil {
		return dtos.BadReqRespMsgCode[*UploadDto](er.Error(), resp_const.InvalidFileType), err
	}
	return dtos.SuccessS(&UploadDto{
		FileType: contentType,
		Hash:     hashString,
		Ext:      MimeToExtension(contentType),
	}, 1), nil
}

type FileType string

const (
	FImageJpg  = FileType("image/jpg")
	FImageJpeg = FileType("image/jpeg")
	FImagePng  = FileType("image/png")
	FImageGif  = FileType("image/gif")
	FImageWebp = FileType("image/webp")
	FilePdf    = FileType("application/pdf")
)

// DetectFileType validates the image content type to check image integrity and validity
func DetectFileType(img multipart.File) (string, error) {

	buffer := make([]byte, 512)

	n, err := img.Read(buffer)
	if err != nil {
		if n < 0 {
			return "", err
		}
		log.Logger = log.With().Caller().Logger()
		log.Err(err).Msg("img type err")
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
	//switch contentType {
	//case jpg, jpge, png, gif:
	//	return contentType, nil
	//default:
	//	return contentType, errors.New("Unsupported CoverImage Format " + contentType + ". Only jpg, jpeg, png and gif is supported")
	//}

}

func MimeToExtension(mimeType string) string {
	mimeMap := map[string]string{
		// Images
		"image/png":     ".png",
		"image/jpeg":    ".jpg",
		"image/jpg":     ".jpg",
		"image/gif":     ".gif",
		"image/webp":    ".webp",
		"image/bmp":     ".bmp",
		"image/x-icon":  ".ico",
		"image/tiff":    ".tiff",
		"image/svg+xml": ".svg",

		// Video
		"video/mp4":        ".mp4",
		"video/mpeg":       ".mpeg",
		"video/ogg":        ".ogv",
		"video/webm":       ".webm",
		"video/x-msvideo":  ".avi",
		"video/x-matroska": ".mkv",

		// Audio
		"audio/mpeg":  ".mp3",
		"audio/ogg":   ".ogg",
		"audio/wav":   ".wav",
		"audio/webm":  ".weba",
		"audio/aac":   ".aac",
		"audio/flac":  ".flac",
		"audio/x-wav": ".wav",

		// Documents
		"text/plain":         ".txt",
		"text/html":          ".html",
		"text/css":           ".css",
		"text/csv":           ".csv",
		"application/json":   ".json",
		"application/xml":    ".xml",
		"application/pdf":    ".pdf",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
		"application/vnd.ms-powerpoint":                                             ".ppt",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",

		// Archives
		"application/zip":              ".zip",
		"application/x-rar-compressed": ".rar",
		"application/x-tar":            ".tar",
		"application/gzip":             ".gz",
		"application/x-7z-compressed":  ".7z",

		// Executables & Code
		"application/x-msdownload":  ".exe",
		"application/x-sh":          ".sh",
		"application/x-python-code": ".py",
		"application/javascript":    ".js",
		"text/x-java-source":        ".java",
		"text/x-c":                  ".c",
		"text/x-c++":                ".cpp",
		"text/x-go":                 ".go",
	}

	// Lookup the extension
	if ext, found := mimeMap[mimeType]; found {
		return ext
	}

	return "" // Fallback: Unknown extension
}
