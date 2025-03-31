package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/google/uuid"
)

type S3Uploader interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, bucketName string) (string, error)
}

type S3Service struct {
	client *s3.Client
}

func NewService(client *s3.Client) *S3Service {
	return &S3Service{
		client: client,
	}
}

func (s S3Service) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, bucketName string) (string, error) {
	fileName := uuid.New().String() + "_" + header.Filename

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		return "", fmt.Errorf("Error reading file: %w", err)
	}

	uploadInput := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(fileName),
		Body:          bytes.NewReader(fileBytes),
		ContentLength: aws.Int64(header.Size),
		ContentType:   aws.String(header.Header.Get("Content-Type")),
	}

	_, err = s.client.PutObject(ctx, uploadInput)
	if err != nil {
		return "", fmt.Errorf("Error while uploading to S3: %w", err)
	}

	return fileName, nil
}
