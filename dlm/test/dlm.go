package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dlm"
	"github.com/aws/aws-sdk-go/service/dlm/dlmiface"
)

type MockDlm struct {
	dlmiface.DLMAPI
	Payload map[string]string // Store expected return values
	Err     error
}

func (d *MockDlm) CreateLifecyclePolicy(i *dlm.CreateLifecyclePolicyInput) (*dlm.CreateLifecyclePolicyOutput, error) {
	if d.Err != nil {
		return nil, d.Err
	}

	return &dlm.CreateLifecyclePolicyOutput{PolicyId: aws.String("test-id")}, nil
}

func (d *MockDlm) DeleteLifecyclePolicy(i *dlm.DeleteLifecyclePolicyInput) (*dlm.DeleteLifecyclePolicyOutput, error) {
	if d.Err != nil {
		return nil, d.Err
	}

	return nil, nil
}

func (d *MockDlm) UpdateLifecyclePolicy(i *dlm.UpdateLifecyclePolicyInput) (*dlm.UpdateLifecyclePolicyOutput, error) {
	if d.Err != nil {
		return nil, d.Err
	}

	return nil, nil
}
