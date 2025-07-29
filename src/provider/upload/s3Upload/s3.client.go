package s3Client

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/birukbelay/gocmn/src/config"
	"github.com/birukbelay/gocmn/src/logger"
	"github.com/birukbelay/gocmn/src/provider/upload"
)

type FileInterface interface {
	SaveFileWithName(data []byte, fileName string) (upload.UploadDto, error)
	SaveFileFromReader(data io.Reader, fileName string) (upload.UploadDto, error)

	DeleteFileWithName(fileName string) error
	GetFileAsBytes(fileName string) ([]byte, error)
	//Functions related to multipart

	InitiateMultipart(fileName string) (string, error)
	AddToMultiPartUpload(data []byte, fileName string, partID int64, uploadID string) (upload.UploadDto, error)
	AddToMultiPartFromReader(data io.ReadSeeker, fileName string, partID int64, uploadID string) (upload.UploadDto, error)
	CompleteMultipartUpload(fileName, uploadID string, completedParts []*s3.CompletedPart) (upload.UploadDto, error)
}

type GarageS3Client struct {
	client     *s3.S3
	bucketName string
	webPath    string
	mu         sync.Mutex
	// uploads    map[string]*multipartUpload
}

func NewAllS3FuncClient(config config.S3Config) (FileInterface, error) {
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
	}, nil
}

func (g *GarageS3Client) SaveFileWithName(data []byte, fileName string) (upload.UploadDto, error) {
	uploader := s3manager.NewUploaderWithClient(g.client)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(g.bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return upload.UploadDto{}, err
	}

	return upload.UploadDto{
		Name: fileName,
		Url:  fmt.Sprintf("%s/%s", g.webPath, fileName),
		Path: result.Location,
		Size: int64(len(data)),
	}, nil
}

func (g *GarageS3Client) SaveFileFromReader(data io.Reader, fileName string) (upload.UploadDto, error) {
	uploader := s3manager.NewUploaderWithClient(g.client)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(g.bucketName),
		Key:    aws.String(fileName),
		Body:   data,
	})
	if err != nil {
		return upload.UploadDto{}, err
	}

	fileUrl := fmt.Sprintf("%s/%s", g.webPath, fileName)
	logger.LogTrace("upload id", result.UploadID)
	// fileUrl := fmt.Sprintf("%s/%s/%s", g.client.Endpoint, g.bucketName, fileName)
	return upload.UploadDto{
		Name:     fileName,
		Url:      fileUrl,
		Path:     result.Location,
		UploadId: result.UploadID,
	}, nil
}

func (g *GarageS3Client) GetFileAsBytes(fileName string) ([]byte, error) {
	resp, err := g.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(g.bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

/* ===========================  MultiPart Uploads   ======================
========================================================================*/

func (g *GarageS3Client) InitiateMultipart(fileName string) (string, error) {
	resp, err := g.client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(g.bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return "", err
	}

	return *resp.UploadId, nil
}
func (g *GarageS3Client) AddToMultiPartUpload(data []byte, fileName string, partID int64, uploadID string) (upload.UploadDto, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Upload part
	resp, err := g.client.UploadPart(&s3.UploadPartInput{
		Bucket:     aws.String(g.bucketName),
		Key:        aws.String(fileName),
		PartNumber: aws.Int64(partID),
		UploadId:   aws.String(uploadID),
		Body:       bytes.NewReader(data),
	})
	if err != nil {
		return upload.UploadDto{}, err
	}

	// g.uploads[streamID].parts = append(g.uploads[streamID].parts, &s3.CompletedPart{
	// 	ETag:       resp.ETag,
	// 	PartNumber: aws.Int64(int64(partID)),
	// })

	// fileUrl := fmt.Sprintf("%s/%s/%s", g.client.Endpoint, g.bucketName, fileName)
	fileUrl := fmt.Sprintf("%s/%s", g.webPath, fileName)

	return upload.UploadDto{
		Name: fileName,
		Url:  fileUrl,
		Path: fileName,
		Size: int64(len(data)),
		ETag: *resp.ETag,
	}, nil
}

func (g *GarageS3Client) AddToMultiPartFromReader(data io.ReadSeeker, fileName string, partID int64, uploadId string) (upload.UploadDto, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Upload part
	resp, err := g.client.UploadPart(&s3.UploadPartInput{
		Bucket:     aws.String(g.bucketName),
		Key:        aws.String(fileName),
		PartNumber: aws.Int64(partID),
		UploadId:   aws.String(uploadId),
		Body:       data,
	})
	if err != nil {
		return upload.UploadDto{}, err
	}

	// g.uploads[uploadId].parts = append(g.uploads[uploadId].parts, &s3.CompletedPart{
	// 	ETag:       resp.ETag,
	// 	PartNumber: aws.Int64(int64(partID)),
	// })

	// fileUrl := fmt.Sprintf("%s/%s/%s", g.client.Endpoint, g.bucketName, fileName)
	fileUrl := fmt.Sprintf("%s/%s", g.webPath, fileName)

	return upload.UploadDto{
		Name: fileName,
		Url:  fileUrl,
		Path: fileName,
		ETag: *resp.ETag,
	}, nil
}

func (g *GarageS3Client) CompleteMultipartUpload(fileName, uploadID string, completedParts []*s3.CompletedPart) (upload.UploadDto, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Complete the multipart upload
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(g.bucketName),
		Key:      aws.String(fileName),
		UploadId: aws.String(uploadID),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}

	_, err := g.client.CompleteMultipartUpload(completeInput)
	if err != nil {
		return upload.UploadDto{}, err
	}

	fileUrl := fmt.Sprintf("%s/%s/%s", g.client.Endpoint, g.bucketName, fileName)

	return upload.UploadDto{
		Name: fileName,
		Url:  fileUrl,
		Path: fmt.Sprintf("%s/%s", g.webPath, fileName),
	}, nil
}
