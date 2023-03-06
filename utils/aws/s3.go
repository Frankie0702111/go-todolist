package aws

import (
	"context"
	"go-todolist/utils/log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func InitS3() *s3.Client {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Panic("Failed to load env file")
	}

	region := os.Getenv("AWS_REGION")
	access := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(access, secret, "")),
		config.WithRegion(region),
	)
	if err != nil {
		log.Error("InitS3 Error: " + err.Error())
		return nil
	}

	return s3.NewFromConfig(cfg)
}
