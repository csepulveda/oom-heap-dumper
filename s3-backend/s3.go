package s3

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func awsSession() (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load AWS configuration: %v", err)
	}
	return cfg, nil
}

// Upload uploads a file to the specified S3 bucket and returns the file's URL.
// It takes the filepath of the file to upload and uses the filename as the S3 object key.
func Upload(bucket, f string) (string, error) {
	// Create an AWS session
	cfg, err := awsSession()
	if err != nil {
		return "", err
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Open the file to upload
	file, err := os.Open(f)
	if err != nil {
		return "", fmt.Errorf("unable to open file %s: %v", f, err)
	}
	defer file.Close()

	// Get the file name to use as the object key in S3
	filename := filepath.Base(f)

	// Upload the file to S3
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload file to S3: %v", err)
	}

	// Construct the URL of the uploaded file
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, filename)
	return url, nil
}
