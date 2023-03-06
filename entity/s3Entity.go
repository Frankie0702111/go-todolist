package entity

import (
	"context"
	"errors"
	"go-todolist/utils/log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/joho/godotenv"
)

type S3Entity interface {
	FileUpload(file *multipart.FileHeader, uuidV4 string) (*manager.UploadOutput, error)
	FileRemove(file string, uuidV4 string) error
}

type s3Connection struct {
	connection *s3.Client
}

func NewS3Entity(db *s3.Client) S3Entity {
	return &s3Connection{
		connection: db,
	}
}

func (db *s3Connection) FileUpload(file *multipart.FileHeader, uuidV4 string) (*manager.UploadOutput, error) {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Panic("Failed to load env file")
	}
	bucker := os.Getenv("AWS_BUCKET")

	client := db.connection
	if client == nil {
		return nil, errors.New("Invalid credential.")
	}

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		// 10 MiB
		u.PartSize = 10 * 1024 * 1024
	})

	f, openErr := file.Open()
	if openErr != nil {
		return nil, openErr
	}

	result, resulterr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucker),
		Key:    aws.String(uuidV4 + "/" + file.Filename),
		Body:   f,
		// ACL:    "public-read",
	})

	if resulterr != nil {
		return nil, resulterr
	}

	return result, nil
}

func (db *s3Connection) FileRemove(file string, uuidV4 string) error {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Panic("Failed to load env file")
	}

	bucker := os.Getenv("AWS_BUCKET")
	var objectIds []types.ObjectIdentifier
	client := db.connection
	if client == nil {
		return errors.New("Invalid credential.")
	}

	objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(uuidV4 + "/" + file)})
	_, err := client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucker),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		return err
	}

	return nil
}
