package upload

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/oklog/ulid/v2"

	file_const "github.com/birukbelay/gocmn/src/consts"
	"github.com/birukbelay/gocmn/src/util"
)

// CreateUploadPath given a file name and extension, trims filename and
func CreateUploadPath(fileName, ext string) (filesName, fullFilePath string, er error) {
	//what happnes if file is
	newFileName := CreateCleanFileNameWithExt(fileName, ext)
	fullPath, err := CreateCleanUploadPath(newFileName)
	if err != nil {
		return "", "", err
	}
	return newFileName, fullPath, nil
}

// CreateCleanFileNameWithExt  cleansup the fiven filename
func CreateCleanFileNameWithExt(fileName, extension string) string {
	//if the filename is empty just return ulid with extension
	if fileName == "" {
		return ulid.Make().String() + extension
	}
	//remove the extension, if it exists
	filenameWithoutExt := strings.TrimSuffix(filepath.Base(fileName), extension)
	//take the first 15 characters, 15 here is pointless, slug only uses first 7
	truncated := filenameWithoutExt[:min(len(filenameWithoutExt), 15)]
	//create a slug, from the first 15 chars
	return util.CreateSlug(truncated) + extension
}

// CreateCleanUploadPath given a file name, just returns the upload path
func CreateCleanUploadPath(fileName string) (filePath string, er error) {
	//cmn.LogTrace("the FileName is", fileName)
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(wd, file_const.FileUploadPath, fileName)
	return fullPath, nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
