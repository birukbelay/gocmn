package diskUpload

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	conf "github.com/birukbelay/gocmn/src/config"
	"github.com/birukbelay/gocmn/src/consts"
	resp "github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/provider/upload"
	"github.com/birukbelay/gocmn/src/resp_const"
)

// SaveFileOnDisk uploads the form file to specific dst. same as Gin
func SaveFileOnDisk(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
func SaveFileToDisk(src multipart.File, dst string) error {

	if err := os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
func DeleteFileWithNameL(fileName string) error {
	// delete local
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filesPath := filepath.Join(wd, consts.FileUploadPath, fileName)
	logger.LogTrace("filePath To remove", filesPath)
	if err := os.Remove(filesPath); err != nil {
		//cmn.LogTrace("Deleting File Error", err)
		log.Err(err).Str("PATH_TO_REMOVE", filesPath).Msg("SingleImage Remove Error")
		//loger.CreateFileLoger(file_const.ImageErrFileName)
		return err
	}
	return nil
}

type LocalUpload struct {
	// Env *conf.EnvConfig
}

func NewDidkUploader(conf *conf.S3Config) upload.FileUploadInterface {
	return LocalUpload{}
}

func (u LocalUpload) UploadFileHeader(file *multipart.FileHeader, opt *upload.Opts) (resp.GResp[*upload.UploadDto], error) {
	//returns file with, fileSize, hash & contentType
	response, err := upload.ValidateFile(file)
	if err != nil {
		return response, err
	}
	//TODO: compress file Here Before Save, & shift size & hash Calculation to here, may be hash should not be moved

	file_name := file.Filename
	if opt != nil {
		if opt.FileName != nil {
			file_name = *opt.FileName
		}
	}
	fileName, filePath, err := upload.CreateUploadPath(file_name, response.Body.Ext)
	if err != nil {
		return resp.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}

	if err := SaveFileOnDisk(file, filePath); err != nil {
		return resp.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.UpdateSuccess), err
	}
	response.Body.Name = fileName
	response.Body.Path = filePath
	response.Body.Url = consts.FileServingUrl + fileName
	return response, nil
}

func (u LocalUpload) DeleteFileWithName(fileName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filesPath := filepath.Join(wd, consts.FileUploadPath, fileName)
	//cmn.LogTrace("filePath To remove", filesPath)
	if err := os.Remove(filesPath); err != nil {
		//cmn.LogTrace("Deleting File Error", err)
		log.Err(err).Str("PATH_TO_REMOVE", filesPath).Msg("SingleImage Remove Error")
		return err
	}
	return nil
}
