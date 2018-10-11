package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dlm"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/liangrog/adlm-helper/dlm/policy"
)

const (
	msgPrefix = "[ADLM-HELPER-INFO]"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) error {
	// Convert context
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		log.Println("Failed to convert context to lambdacontext")
	}

	p := new(policy.Policy)

	// Initiate aws client
	sess := session.Must(session.NewSession())
	clients := &policy.AwsClients{
		S3Downloader: s3manager.NewDownloader(sess),
		Dynamodb:     dynamodb.New(sess),
		Dlm:          dlm.New(sess),
	}

	p.SetClients(clients)

	errCount := 0
	for _, record := range s3Event.Records {
		// Ignore if event triggered is by directory lifecycle
		if isDir(record.S3.Object.Key) {
			log.Println(fmt.Sprintf("%s Ignoring event triggered by %s because it is not a file", msgPrefix, record.S3.Object.Key))
			continue
		}

		p.SetPolicy(record, *lc)
		err := p.Dispatch().Execute()
		if err != nil {
			// Log event for debugging
			log.Println(fmt.Sprintf("%s %v", msgPrefix, record))

			errCount++

			// Logging event for debugging
			log.Println(fmt.Sprintf("%s %v", msgPrefix, record))

			if awsErr, ok := err.(awserr.Error); ok {
				log.Println("Error:", awsErr.Error())
			} else {
				log.Println(err)
			}
		} else {
			log.Println(fmt.Sprintf("%s Successfully processed event triggered by %s", msgPrefix, record.S3.Object.Key))
		}
	}

	if errCount > 0 {
		return errors.New(fmt.Sprintf("%s Failed processed total record: %d", msgPrefix, errCount))
	}

	return nil
}

// If it's a S3 directory
func isDir(s string) bool {
	return strings.HasSuffix(s, "/")
}
