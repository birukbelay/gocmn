package cloudinaryServ

import (
	"context"
	"mime/multipart"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/rs/zerolog/log"

	conf "github.com/birukbelay/gocmn/src/config"
	resp "github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/provider/upload"
	"github.com/birukbelay/gocmn/src/resp_const"
)

// CloudinaryUpload struct to hold Cloudinary client and config
type CloudinaryUpload struct {
	Env   *conf.CloudinaryConfig
	Cloud *cloudinary.Cloudinary
}

func NewCloudinaryUploader(conf *conf.CloudinaryConfig) upload.FileUploadInterface {
	// Initialize Cloudinary client with credentials from EnvConfig
	cld, err := cloudinary.NewFromParams(conf.CloudinaryCloudName, conf.CloudinaryAPIKey, conf.CloudinaryAPISecret)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Cloudinary")
	}
	cld.Config.URL.Secure = true
	return &CloudinaryUpload{
		Env:   conf,
		Cloud: cld,
	}
}

// UploadSingleFile uploads a file to Cloudinary with a generated name
func (u *CloudinaryUpload) UploadSingleFile(file *multipart.FileHeader) (resp.GResp[upload.UploadDto], error) {
	// Validate file (reusing your existing function)
	response, err := upload.ValidateFile(file)
	if err != nil {
		return response, err
	}

	// Open the file
	f, err := file.Open()
	if err != nil {
		return resp.BadReqRespMsgCode[upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}
	defer f.Close()

	// Upload to Cloudinary
	ctx := context.Background()
	uploadResult, err := u.Cloud.Upload.Upload(ctx, f, uploader.UploadParams{
		Folder: u.Env.CloudinaryFolder, // Optional: e.g., "uploads" from config
	})
	if err != nil {
		return resp.BadReqRespMsgCode[upload.UploadDto](err.Error(), resp_const.UpdateSuccess), err
	}

	// Populate response
	response.Body.Name = filepath.Base(uploadResult.PublicID) // e.g., "random123"
	response.Body.Path = uploadResult.PublicID                // Full Cloudinary path, e.g., "uploads/random123"
	response.Body.Url = uploadResult.SecureURL                // HTTPS URL for access
	response.Body.Size = int64(uploadResult.Bytes)            // Update size from Cloudinary
	return response, nil
}

// UploadSingleFile uploads a file to Cloudinary with a generated name
func (u *CloudinaryUpload) UploadFile(file multipart.File) (resp.GResp[upload.UploadDto], error) {
	// Validate file (reusing your existing function)
	response, err := upload.ValidateJustFile(file)
	if err != nil {
		return response, err
	}
	// Upload to Cloudinary
	ctx := context.Background()
	uploadResult, err := u.Cloud.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: u.Env.CloudinaryFolder, // Optional: e.g., "uploads" from config
	})
	if err != nil {
		return resp.BadReqRespMsgCode[upload.UploadDto](err.Error(), resp_const.UpdateSuccess), err
	}
	// Populate response
	response.Body.Name = filepath.Base(uploadResult.PublicID) // e.g., "random123"
	response.Body.Path = uploadResult.PublicID                // Full Cloudinary path, e.g., "uploads/random123"
	response.Body.Url = uploadResult.SecureURL                // HTTPS URL for access
	response.Body.Size = int64(uploadResult.Bytes)            // Update size from Cloudinary
	return response, nil
}

// UploadWithGivenName uploads a file to Cloudinary with a specified name
func (u *CloudinaryUpload) UploadWithGivenName(file *multipart.FileHeader, name string) (resp.GResp[upload.UploadDto], error) {
	// Validate file
	response, err := upload.ValidateFile(file)
	if err != nil {
		return response, err
	}

	fileName := name + response.Body.Ext

	// Open the file
	f, err := file.Open()
	if err != nil {
		return resp.BadReqRespMsgCode[upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}
	defer f.Close()

	// Upload to Cloudinary with custom PublicID
	ctx := context.Background()
	uploadResult, err := u.Cloud.Upload.Upload(ctx, f, uploader.UploadParams{
		PublicID: fileName,               // Custom name without folder prefix
		Folder:   u.Env.CloudinaryFolder, // Optional: e.g., "uploads"
	})
	if err != nil {
		return resp.BadReqRespMsgCode[upload.UploadDto](err.Error(), resp_const.UpdateSuccess), err
	}

	// Populate response
	response.Body.Name = filepath.Base(uploadResult.PublicID) // e.g., "customname.jpg"
	response.Body.Path = uploadResult.PublicID                // e.g., "uploads/customname.jpg"
	response.Body.Url = uploadResult.SecureURL                // HTTPS URL
	response.Body.Size = int64(uploadResult.Bytes)            // Size from Cloudinary
	return response, nil
}

// DeleteFileWithName deletes a file from Cloudinary by its name (PublicID)
func (u *CloudinaryUpload) DeleteFileWithName(fileName string) error {
	ctx := context.Background()
	// Assume fileName is the PublicID (e.g., "uploads/myfile.jpg" or just "myfile.jpg")
	_, err := u.Cloud.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: fileName,
	})
	if err != nil {
		log.Err(err).Str("PUBLIC_ID", fileName).Msg("Cloudinary Delete Error")
		return err
	}
	return nil
}
