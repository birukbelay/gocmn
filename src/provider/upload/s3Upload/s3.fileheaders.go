package s3Client

import (
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/birukbelay/gocmn/src/config"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/provider/upload"
	"github.com/birukbelay/gocmn/src/resp_const"
)

func NewLocS3FuncClient(config config.S3Config) (upload.FileUploadInterface, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(config.Endpoint),
		Region:           aws.String(config.Region),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess)
	return &GarageS3Client{
		client:     s3Client,
		bucketName: config.BucketName,
		webPath:    config.S3WebEndpoint,
		// uploads:    make(map[string]*multipartUpload),
	}, nil
}

func (g *GarageS3Client) UploadFileHeader(file *multipart.FileHeader, opt *upload.Opts) (dtos.GResp[*upload.UploadDto], error) {
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
	fileName := upload.CreateCleanFileNameWithExt(file_name, response.Body.Ext)

	src, err := file.Open()
	if err != nil {
		return dtos.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}
	defer src.Close()
	uploadRes, err := g.SaveFileFromReader(src, fileName)
	if err != nil {
		return dtos.BadReqRespMsgCode[*upload.UploadDto](err.Error(), resp_const.CantReadFileType), err
	}

	response.Body.Name = uploadRes.Name
	response.Body.Url = uploadRes.Url

	return response, nil
}

func (g *GarageS3Client) DeleteFileWithName(fileName string) error {
	_, err := g.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(g.bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return err
	}

	return nil
}
