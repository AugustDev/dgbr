package awsx

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Config - contains details required to interact with AWS S3 bucket
type Config struct {
	Region       string
	Bucket       string
	IAMAccessKey string
	IAMSecretKey string
}

// CreateSession - establishes AWS session and checks credentials
func (conf *Config) CreateSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      &conf.Region,
		Credentials: credentials.NewStaticCredentials(conf.IAMAccessKey, conf.IAMSecretKey, ""),
	})
	if err != nil {
		return nil, errors.New("could not create AWS session, check your credentials")
	}

	return sess, nil
}

// UploadToS3 - uploads backup filo to AWS S3 bucket
func (conf *Config) UploadToS3(filepath string) error {

	sess, err := conf.CreateSession()
	if err != nil {
		return err
	}

	// initializes S3 uploader
	uploader := s3manager.NewUploader(sess)

	// reading file contents
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file %q with error: %v", filepath, err)
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: &conf.Bucket,
		Key:    &filepath,
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload the backup with error: %v", err)
	}

	return nil
}

// DownloadFromS3 - downloads backup file from S3 bucket
func (conf *Config) DownloadFromS3(key string) (string, error) {

	sess, err := conf.CreateSession()
	if err != nil {
		return "", err
	}

	downloader := s3manager.NewDownloader(sess)
	file, err := os.Create(key)
	if err != nil {
		return "", fmt.Errorf("could not write to file with error: %v", err)
	}
	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: &conf.Bucket,
			Key:    &key,
		})
	if err != nil {
		return "", fmt.Errorf("could not download file %q, with error: %v", key, err)
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	return file.Name(), nil
}
