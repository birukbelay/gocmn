package firebaseUpload

import (
	"context"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	firestorage "firebase.google.com/go/v4/storage"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/provider/upload"
	"github.com/birukbelay/gocmn/src/resp_const"
)

type FirebaseUpload struct {
	App        *firebase.App
	BucketName string
	client     *firestorage.Client
	bucket     *storage.BucketHandle
}

func NewFirebaseUpload(app *firebase.App, bucketName string) (upload.FileUploadInterface, error) {
	ctx := context.Background()
	client, err := app.Storage(ctx)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return &FirebaseUpload{
		App:        app,
		BucketName: bucketName,
		client:     client,
		bucket:     bucket,
	}, nil
}

func (f *FirebaseUpload) UploadFile(src multipart.File, opt *upload.Opts) (dtos.GResp[*upload.UploadDto], error) {
	// Re-verify file type and get info using your existing utility
	response, err := upload.GetFileInfo(src)
	if err != nil {
		return response, err
	}

	fileName := ""
	if opt != nil && opt.FileName != nil {
		fileName = *opt.FileName
	}

	cleanFileName := upload.CreateCleanFileNameWithExt(fileName, response.Body.Ext)

	// Since we already read the file to calculate hash/type in GetFileInfo, we must seek back to start
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return dtos.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}

	return f.uploadToFirebase(src, cleanFileName, response, response.Body.FileType)
}

func (f *FirebaseUpload) UploadFileHeader(file *multipart.FileHeader, opt *upload.Opts) (dtos.GResp[*upload.UploadDto], error) {
	// Validate file (reusing your existing function)
	response, err := upload.ValidateFile(file)
	if err != nil {
		return response, err
	}

	fileName := file.Filename
	if opt != nil && opt.FileName != nil {
		fileName = *opt.FileName
	}

	cleanFileName := upload.CreateCleanFileNameWithExt(fileName, response.Body.Ext)

	// Open the file
	src, err := file.Open()
	if err != nil {
		return dtos.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}
	defer src.Close()

	return f.uploadToFirebase(src, cleanFileName, response, response.Body.FileType)
}

func (f *FirebaseUpload) uploadToFirebase(src io.Reader, fileName string, response dtos.GResp[*upload.UploadDto], contentType string) (dtos.GResp[*upload.UploadDto], error) {
	ctx := context.Background()

	obj := f.bucket.Object(fileName)
	writer := obj.NewWriter(ctx)
	writer.ContentType = contentType

	if _, err := io.Copy(writer, src); err != nil {
		writer.Close()
		return dtos.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.UpdateSuccess), err
	}

	if err := writer.Close(); err != nil {
		return dtos.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.UpdateSuccess), err
	}

	// Make the file publicly accessible (optional, usually required for uploads)
	if err := obj.ACL().Set(ctx, "allUsers", "READER"); err != nil {
		// Log error but don't fail upload
	}

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return dtos.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.UpdateSuccess), err
	}

	// Output URL logic: Firebase Storage public URLs look like this
	response.Body.Name = fileName
	response.Body.Path = fileName
	response.Body.Url = attrs.MediaLink // This provides a download URL
	response.Body.Size = attrs.Size

	return response, nil
}

func (f *FirebaseUpload) DeleteFileWithName(fileName string) error {
	ctx := context.Background()

	return f.bucket.Object(fileName).Delete(ctx)
}
