package file

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"github.com/liangrog/adlm-helper/dlm/test"
)

var (
	record = events.S3EventRecord{
		S3: events.S3Entity{
			Bucket: events.S3Bucket{
				Name: "dummy-bucket",
			},
			Object: events.S3Object{
				Key: test.PolicyExampleFileName,
			},
		},
	}
)

func TestS3Downloader(t *testing.T) {
	err := S3Downloader(test.DestTestFile, record.S3.Bucket.Name, record.S3.Object.Key, new(test.MockDownloader))
	assert.NoError(t, err)

	//Clean up test file
	err = test.DeleteFile(test.DestTestFile)
	assert.NoError(t, err)
}

func TestUnmarshalPolicyFromS3(t *testing.T) {
	p, err := UnmarshalPolicyFromS3(record, new(test.MockDownloader))
	assert.NoError(t, err)
	assert.Equal(t, "My Awesome Data Lifecycl Management Daily Snapshot", p.Description, "Policy description not match")
	assert.Equal(t, "VOLUME", p.PolicyDetails.ResourceTypes, "Policy ResourceTypes not match")

	// Clean up test file
	err = test.DeleteFile(test.DestTestFile)
	assert.NoError(t, err)
}
