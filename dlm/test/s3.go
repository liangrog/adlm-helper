package test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

const (
	// This is the relative path to github.com/liangrog/adlm-helper/dlm/file/file_test.go
	// as the test is run from that folder
	policyExampleFileSourcePath = "../../testdata"

	PolicyExampleFileName = "policy_example.yaml"

	cacheDir = "/tmp"
)

var (
	SrcTestFile  = path.Join(policyExampleFileSourcePath, PolicyExampleFileName)
	DestTestFile = path.Join(cacheDir, PolicyExampleFileName)
)

// Mocking Downloader
type MockDownloader struct {
	s3manageriface.DownloaderAPI
}

func (md *MockDownloader) Download(iw io.WriterAt, gi *s3.GetObjectInput, dl ...func(*s3manager.Downloader)) (int64, error) {
	return CopyFile(SrcTestFile, DestTestFile)
}

// Copy files locally
func CopyFile(src, dest string) (int64, error) {
	var nbyte int64
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Printf("Can't read file from %s. Error: %v", src, err)
		return nbyte, err
	}

	nbyte = int64(len(input))

	err = ioutil.WriteFile(dest, input, 0644)
	if err != nil {
		fmt.Printf("Can't write file to %s. Error: %v", dest, err)
		return nbyte, err
	}

	return nbyte, nil
}

// Delete local file
func DeleteFile(src string) error {
	return os.Remove(src)
}
