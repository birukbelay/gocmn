package upload

import (
	"os"
	"path/filepath"
	"strings"

	file_const "github.com/birukbelay/gocmn/src/consts"
	"github.com/birukbelay/gocmn/src/util"
)

func CreateUploadPath(fileName string) (filesName, fullFilePath string, er error) {
	//cmn.LogTrace("the FileName is", fileName)
	newFileName := CreateFileName(fileName)
	wd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	fullPath := filepath.Join(wd, file_const.FileUploadPath, newFileName)
	return newFileName, fullPath, nil
}
func CreateCleanUploadPath(fileName string) (filePath string, er error) {
	//cmn.LogTrace("the FileName is", fileName)
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(wd, file_const.FileUploadPath, fileName)
	return fullPath, nil
}
func CreateFileName(fileName string) string {
	extension := filepath.Ext(fileName)
	filenameWithoutExt := strings.TrimSuffix(filepath.Base(fileName), extension)
	truncated := filenameWithoutExt[:min(len(filenameWithoutExt), 15)]
	return util.CreateSlug(truncated) + extension
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
