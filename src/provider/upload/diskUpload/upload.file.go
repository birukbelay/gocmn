package diskUpload

import (
	"mime/multipart"

	"github.com/birukbelay/gocmn/src/consts"
	resp "github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/provider/upload"
	"github.com/birukbelay/gocmn/src/resp_const"
)

func (u LocalUpload) UploadFile(file multipart.File, opt *upload.Opts) (resp.GResp[*upload.UploadDto], error) {
	//returns file with, fileSize, hash & contentType
	response, err := upload.GetFileInfo(file)
	if err != nil {
		return response, err
	}
	//TODO: compress file Here Before Save, & shift size & hash Calculation to here, may be hash should not be moved

	file_name := ""
	if opt != nil {
		if opt.FileName != nil {
			file_name = *opt.FileName
		}
	}
	fileName, filePath, err := upload.CreateUploadPath(file_name, response.Body.Ext)
	if err != nil {
		return resp.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}

	if err := SaveFileToDisk(file, filePath); err != nil {
		return resp.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.UpdateSuccess), err
	}
	response.Body.Name = fileName
	response.Body.Path = filePath
	response.Body.Url = consts.FileServingUrl + fileName
	return response, nil
}
