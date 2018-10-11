package db

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Database abstract
type DB interface {
	FindByKey(string) (*Item, error)
	Create(*Item) error
	Update(*Item) error
	Delete(*Item) error
}

// Database item
type Item struct {
	S3ObjectKey string `json:"s3objectkey"`
	PolicyId    string `json:"policyid"`
	RequestId   string `json:"requestid"`
	CreatedAt   string `json:"createdat"`
	UpdatedAt   string `json:"updatedat"`
}

// Database factory
func GetConn(i interface{}) DB {
	switch v := i.(type) {
	case dynamodbiface.DynamoDBAPI:
		return &Dynamo{
			client: v,
		}
	}

	return nil
}
