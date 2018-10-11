package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// DynamoDB table name. Must be the same as template.yaml
var (
	tableName = "adlm-helper"
)

// Table primary key
type ItemKey struct {
	S3ObjectKey string `json:"s3objectkey"`
}

// fields needed for update
type ItemUpdate struct {
	RequestId string `json:":r"`
	UpdatedAt string `json:":u"`
}

type Dynamo struct {
	client dynamodbiface.DynamoDBAPI
}

// Search by key
func (d *Dynamo) FindByKey(k string) (*Item, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":vk": {
				S: aws.String(k),
			},
		},
		KeyConditionExpression: aws.String("s3objectkey = :vk"),
		TableName:              aws.String(tableName),
	}

	result, err := d.client.Query(input)
	if err != nil {
		return nil, err
	}

	if result.Items != nil && len(result.Items) > 0 {
		item := new(Item)
		if err = dynamodbattribute.UnmarshalMap(result.Items[0], item); err != nil {
			return nil, err
		}

		return item, nil
	}

	return nil, nil
}

// Create a record
func (d *Dynamo) Create(i *Item) error {
	item, err := dynamodbattribute.MarshalMap(i)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = d.client.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// Update a record
func (d *Dynamo) Update(i *Item) error {
	key, err := dynamodbattribute.MarshalMap(ItemKey{
		S3ObjectKey: i.S3ObjectKey,
	})

	if err != nil {
		return err
	}

	update, err := dynamodbattribute.MarshalMap(ItemUpdate{
		RequestId: i.RequestId,
		UpdatedAt: i.UpdatedAt,
	})

	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		Key: key,
		ExpressionAttributeNames: map[string]*string{
			"#RI": aws.String("requestid"),
			"#UA": aws.String("updatedat"),
		},
		ExpressionAttributeValues: update,
		TableName:                 aws.String(tableName),
		UpdateExpression:          aws.String("SET #RI = :r, #UA = :u"),
	}

	_, err = d.client.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}

// Delete a record
func (d *Dynamo) Delete(i *Item) error {
	key, err := dynamodbattribute.MarshalMap(ItemKey{
		S3ObjectKey: i.S3ObjectKey,
	})

	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(tableName),
	}

	_, err = d.client.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}
