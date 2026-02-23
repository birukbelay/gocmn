package cloudinaryServ

import (
	"context"
	"mime/multipart"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	resp "github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/provider/upload"
	"github.com/birukbelay/gocmn/src/resp_const"
)

// UploadFileHeader uploads a file to Cloudinary with a generated name
func (u *CloudinaryUpload) UploadFileHeader(file *multipart.FileHeader, opt *upload.Opts) (resp.GResp[*upload.UploadDto], error) {
	// Validate file (reusing your existing function)
	response, err := upload.ValidateFile(file)
	if err != nil {
		return response, err
	}

	file_name := file.Filename
	if opt != nil {
		if opt.FileName != nil {
			file_name = *opt.FileName
		}
	}
	cloudinaryParams := uploader.UploadParams{
		Folder:       u.Env.CloudinaryFolder, // Optional: e.g., "uploads" from config
		PublicID:     upload.CreateCleanFileNameWithExt(file_name, response.Body.Ext),
		ResourceType: "auto",
		Type:         "upload",

		AccessControl: api.AccessControl{
			{AccessType: "anonymous"},
		},
	}
	// Open the file
	f, err := file.Open()
	if err != nil {
		return resp.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}
	defer f.Close()

	// Upload to Cloudinary
	ctx := context.Background()
	uploadResult, err := u.Cloud.Upload.Upload(ctx, f, cloudinaryParams)
	if err != nil {
		return resp.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.UpdateSuccess), err
	}
	// signedURL, err := u.Cloud..URL.Signed(uploadResult.PublicID, &url.Configuration{
	// 	ResourceType: uploadResult.ResourceType,
	// 	Secure:       true,
	// 	SignURL:      true,
	// })

	// Populate response
	response.Body.Name = filepath.Base(uploadResult.PublicID) // e.g., "random123"
	response.Body.Path = uploadResult.PublicID                // Full Cloudinary path, e.g., "uploads/random123"
	response.Body.Url = uploadResult.SecureURL                // HTTPS URL for access
	response.Body.Size = int64(uploadResult.Bytes)            // Update size from Cloudinary
	return response, nil
}
