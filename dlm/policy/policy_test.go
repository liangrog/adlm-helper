package policy

import (
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/service/dlm"
	"github.com/stretchr/testify/assert"

	"github.com/liangrog/adlm-helper/dlm/db"
	"github.com/liangrog/adlm-helper/dlm/test"
)

// Mock event record
var record = events.S3EventRecord{
	EventTime: time.Now(),
	S3: events.S3Entity{
		Bucket: events.S3Bucket{
			Name: "dummy-bucket",
		},
		Object: events.S3Object{
			Key: test.PolicyExampleFileName,
		},
	},
}

// Mock event context
var context = lambdacontext.LambdaContext{
	AwsRequestID: "abcde-12345",
}

// Mock database item
var item = &db.Item{
	S3ObjectKey: "test",
	PolicyId:    "test-id",
	RequestId:   "request-id",
	CreatedAt:   "1900-00-00 00:00:00",
	UpdatedAt:   "1900-00-00 00:00:00",
}

// Mock clients
func GetClients(hasDbItem bool) *AwsClients {
	var dynamo *test.MockDynamoDB
	if hasDbItem {
		dynamo = &test.MockDynamoDB{
			Payload: map[string]string{
				"found": "yes",
			},
		}
	} else {
		dynamo = new(test.MockDynamoDB)
	}

	return &AwsClients{
		S3Downloader: new(test.MockDownloader),
		Dynamodb:     dynamo,
		Dlm:          new(test.MockDlm),
	}
}

// Mock policy
func GetPolicy(hasDbItem bool) *Policy {
	p := new(Policy)
	p.SetClients(GetClients(hasDbItem))
	p.SetPolicy(record, context)
	return p
}

func GetDeleterProcessor() Processor {
	record.EventName = "ObjectRemoved:DeleteMarkerCreated"
	p := GetPolicy(true)
	return p.Dispatch()
}

func GetUpserterProcessor(hasDbItem bool) Processor {
	record.EventName = "ObjectCreated:Put"
	p := GetPolicy(hasDbItem)
	return p.Dispatch()
}

func TestPolicyFacadeDeleter(t *testing.T) {
	proc := GetDeleterProcessor()
	assert.IsType(t, Deleter{}, proc, "Deleter type doesn't match")
}

func TestPolicyFacadeUpserter(t *testing.T) {
	proc := GetUpserterProcessor(false)
	assert.IsType(t, Upserter{}, proc, "Upserter type doesn't match")
}

func TestUpserterHydrateCreate(t *testing.T) {
	proc := GetUpserterProcessor(false)

	upserter, ok := proc.(Upserter)
	assert.True(t, ok)

	input, err := upserter.hydrate()
	assert.NoError(t, err)
	assert.IsType(t, &dlm.CreateLifecyclePolicyInput{}, input)
}

func TestUpserterHydrateUpdate(t *testing.T) {
	proc := GetUpserterProcessor(true)

	upserter, ok := proc.(Upserter)
	assert.True(t, ok)

	input, err := upserter.hydrate()
	assert.NoError(t, err)
	assert.IsType(t, &dlm.UpdateLifecyclePolicyInput{}, input)
}

func TestCreatPolicy(t *testing.T) {
	proc := GetUpserterProcessor(false)

	upserter, ok := proc.(Upserter)
	assert.True(t, ok)

	err := upserter.CreatePolicy()
	assert.NoError(t, err)
}

func TestUpdatePolicy(t *testing.T) {
	proc := GetUpserterProcessor(true)

	upserter, ok := proc.(Upserter)
	assert.True(t, ok)

	err := upserter.UpdatePolicy()
	assert.NoError(t, err)
}

func TestDeletePolicy(t *testing.T) {
	proc := GetDeleterProcessor()

	deleter, ok := proc.(Deleter)
	assert.True(t, ok)

	err := deleter.DeletePolicy()
	assert.NoError(t, err)
}
