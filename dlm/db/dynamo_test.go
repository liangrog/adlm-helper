package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	test "github.com/liangrog/adlm-helper/dlm/test"
)

var it = &Item{
	S3ObjectKey: "test",
	PolicyId:    "test-id",
	RequestId:   "request-id",
	CreatedAt:   "1900-00-00 00:00:00",
	UpdatedAt:   "1900-00-00 00:00:00",
}

func TestFindByKeyItemFound(t *testing.T) {
	dy := &Dynamo{
		client: &test.MockDynamoDB{
			Payload: map[string]string{
				"found": "yes",
			},
		},
	}

	i, err := dy.FindByKey("dummpy-key")
	assert.NoError(t, err)
	assert.Equal(t, "found", i.S3ObjectKey, "s3 key doesn't match")
}

func TestFindByKeyNoItemFound(t *testing.T) {
	dy := &Dynamo{
		client: &test.MockDynamoDB{},
	}

	i, err := dy.FindByKey("dummpy-key")
	assert.NoError(t, err)
	assert.Nil(t, i)
}

func TestFindByKeyError(t *testing.T) {
	dy := &Dynamo{
		client: &test.MockDynamoDB{
			Err: errors.New("error"),
		},
	}

	_, err := dy.FindByKey("dummpy-key")
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {
	dy := &Dynamo{
		client: &test.MockDynamoDB{},
	}

	err := dy.Create(it)
	assert.NoError(t, err)
}

func TestCreateError(t *testing.T) {
	dy := &Dynamo{
		client: &test.MockDynamoDB{
			Err: errors.New("error"),
		},
	}

	err := dy.Create(it)
	assert.Error(t, err)
}

func TestUpdate(t *testing.T) {
	dy := &Dynamo{
		client: &test.MockDynamoDB{},
	}

	err := dy.Update(it)
	assert.NoError(t, err)

}

func TestUpdateError(t *testing.T) {
	dy := &Dynamo{
		client: &test.MockDynamoDB{
			Err: errors.New("error"),
		},
	}

	err := dy.Update(it)
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	dy := &Dynamo{
		client: &test.MockDynamoDB{},
	}

	err := dy.Delete(it)
	assert.NoError(t, err)

}

func TestDeleteError(t *testing.T) {
	dy := &Dynamo{
		client: &test.MockDynamoDB{
			Err: errors.New("error"),
		},
	}

	err := dy.Delete(it)
	assert.Error(t, err)
}
