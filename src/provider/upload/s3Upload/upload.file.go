package s3Client

import (
	"mime/multipart"

	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/provider/upload"
	"github.com/birukbelay/gocmn/src/resp_const"
)

func (g *GarageS3Client) UploadFile(src multipart.File, opt *upload.Opts) (dtos.GResp[*upload.UploadDto], error) {
	response, err := upload.GetFileInfo(src)
	if err != nil {
		return response, err
	}

	filename := ""
	if opt != nil {
		if opt.FileName != nil {
			filename = *opt.FileName
		}
	}
	fileName := upload.CreateCleanFileNameWithExt(filename, response.Body.Ext)

	uploadRes, err := g.SaveFileFromReader(src, fileName)
	if err != nil {
		return dtos.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}

	response.Body.Name = uploadRes.Name
	response.Body.Url = uploadRes.Url

	return response, nil
}
