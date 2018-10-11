package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Mocking dynamoDB
type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	Payload map[string]string // Store expected return values
	Err     error
}

func (m MockDynamoDB) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	// Test found item condition
	var output *dynamodb.QueryOutput
	if m.Payload["found"] == "yes" {
		item := map[string]*dynamodb.AttributeValue{
			"S3ObjectKey": {
				S: aws.String("found"),
			},
			"PolicyId": {
				S: aws.String("abcde-12345"),
			},
		}

		output = &dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				item,
			},
		}
	} else {
		output = &dynamodb.QueryOutput{}
	}

	return output, nil
}

func (m MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &dynamodb.PutItemOutput{}, nil
}

func (m MockDynamoDB) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &dynamodb.UpdateItemOutput{}, nil
}

func (m MockDynamoDB) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return &dynamodb.DeleteItemOutput{}, nil
}
