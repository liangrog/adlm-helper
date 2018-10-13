package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"gopkg.in/yaml.v2"
)

const cacheDir = "/tmp"

// Unmarshal yaml file from local directory after downloaded it
func UnmarshalPolicyFromS3(record events.S3EventRecord, downloader s3manageriface.DownloaderAPI) (*Policy, error) {
	localFile := filepath.Join(cacheDir, record.S3.Object.Key)

	// Download file to lambda container temperarily
	if err := S3Downloader(localFile, record.S3.Bucket.Name, record.S3.Object.Key, downloader); err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadFile(localFile)
	if err != nil {
		return nil, err
	}

	p := new(Policy)
	if err = yaml.Unmarshal(raw, p); err != nil {
		return nil, err
	}

	return p, nil
}

// Download file from S3 bucket
func S3Downloader(fileName, bucket, key string, downloader s3manageriface.DownloaderAPI) error {
	// Create a file to write the S3 Object contents to.
	if err := os.MkdirAll(filepath.Dir(fileName), 0777); err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", fileName, err)
	}

	// Write the contents of S3 Object to the file
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}

	return nil
}
