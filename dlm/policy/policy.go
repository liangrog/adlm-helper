package policy

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dlm"
	"github.com/aws/aws-sdk-go/service/dlm/dlmiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"

	"github.com/liangrog/adlm-helper/dlm/db"
	"github.com/liangrog/adlm-helper/dlm/file"
)

// Item from handler given by triggered events
type eventItem struct {
	record  events.S3EventRecord
	context lambdacontext.LambdaContext
	dbItem  *db.Item
}

// AWS services client
type AwsClients struct {
	S3Downloader s3manageriface.DownloaderAPI
	Dynamodb     dynamodbiface.DynamoDBAPI
	Dlm          dlmiface.DLMAPI
}

// Policy entity.
// Facade of the processors
type Policy struct {
	client *AwsClients
	item   *eventItem
	dbconn db.DB
}

// Set AWS client
func (p *Policy) SetClients(c *AwsClients) {
	p.client = c
	// Initiate the database session
	p.dbconn = db.GetConn(p.client.Dynamodb)
}

// Set DB Conn
func (p *Policy) SetDBConn(c db.DB) {
	p.dbconn = db.GetConn(c)
}

// Policy setter
func (p *Policy) SetPolicy(r events.S3EventRecord, c lambdacontext.LambdaContext) error {
	p.item = &eventItem{
		record:  r,
		context: c,
	}

	// Proactive checking if policy exist
	if di, err := p.dbconn.FindByKey(p.item.record.S3.Object.Key); err != nil {
		return err
	} else if di != nil {
		p.item.dbItem = di
	} else {
		p.item.dbItem = nil
	}

	return nil
}

// Decide what to do with a given s3 event record.
// It returns the processor for invoking
func (p *Policy) Dispatch() Processor {
	// If EventName start with ObjectRemoved, it indicates it's a delete event
	if re := regexp.MustCompile(`^ObjectRemoved`); re.MatchString(p.item.record.EventName) {
		return Deleter{
			item:   p.item,
			client: p.client,
			dbconn: p.dbconn,
		}
	}

	// Everything else is create/update event
	return Upserter{
		item:   p.item,
		client: p.client,
		dbconn: p.dbconn,
	}
}

// Strategy pattern
type Processor interface {
	Execute() error
}

// Strategy upsert
// Perform policy create or update
// The event s3 generate doesn't distinguish
// if it's a file update or create. Hence we
// need to relying on checking DB records
type Upserter struct {
	item   *eventItem
	client *AwsClients
	dbconn db.DB
}

func (u Upserter) Execute() error {
	if u.item.dbItem == nil {
		return u.CreatePolicy()
	}

	return u.UpdatePolicy()
}

// Populate the input from records.
// The return value can be create input or update input depends on the event.
func (u Upserter) hydrate() (interface{}, error) {
	// Load policy config from s3
	f, err := file.UnmarshalPolicyFromS3(u.item.record, u.client.S3Downloader)
	if err != nil {
		return nil, err
	}

	// Retain Rule
	retainRule := new(dlm.RetainRule).
		SetCount(f.PolicyDetails.Schedules[0].RetainRule.Count)

	// Create Rule
	createRule := new(dlm.CreateRule).
		SetInterval(f.PolicyDetails.Schedules[0].CreateRule.Interval).
		SetIntervalUnit(f.PolicyDetails.Schedules[0].CreateRule.IntervalUnit).
		SetTimes(f.PolicyDetails.Schedules[0].CreateRule.Times)

	// Schedules
	var tagsToAdd []*dlm.Tag
	for _, t := range f.PolicyDetails.Schedules[0].TagsToAdd {
		tagsToAdd = append(tagsToAdd, &dlm.Tag{Key: aws.String(t.Key), Value: aws.String(t.Value)})
	}

	schedule := new(dlm.Schedule).
		SetName(f.PolicyDetails.Schedules[0].Name).
		SetCreateRule(createRule).
		SetRetainRule(retainRule).
		SetTagsToAdd(tagsToAdd)

	// PolicyDetails
	var targetTags []*dlm.Tag
	for _, t := range f.PolicyDetails.TargetTags {
		targetTags = append(targetTags, &dlm.Tag{Key: aws.String(t.Key), Value: aws.String(t.Value)})
	}

	policyDetails := new(dlm.PolicyDetails).
		SetResourceTypes([]*string{aws.String(f.PolicyDetails.ResourceTypes)}).
		SetSchedules([]*dlm.Schedule{schedule}).
		SetTargetTags(targetTags)

	// If it's update
	if u.item.dbItem != nil {
		return new(dlm.UpdateLifecyclePolicyInput).
				SetDescription(f.Description).
				SetExecutionRoleArn(f.ExecutionRoleArn).
				SetPolicyDetails(policyDetails).
				SetState(f.State).
				SetPolicyId(u.item.dbItem.PolicyId),
			nil
	}

	// If it's create
	return new(dlm.CreateLifecyclePolicyInput).
			SetDescription(f.Description).
			SetExecutionRoleArn(f.ExecutionRoleArn).
			SetPolicyDetails(policyDetails).
			SetState(f.State),
		nil
}

// Create DLM polocy and save the result into database
func (u Upserter) CreatePolicy() error {
	var input *dlm.CreateLifecyclePolicyInput

	i, err := u.hydrate()
	if err != nil {
		return err
	}

	input, ok := i.(*dlm.CreateLifecyclePolicyInput)
	if !ok {
		return errors.New("Failed to cast data into CreateLifecyclePolicyInput")
	}

	// Create policy
	output, err := u.client.Dlm.CreateLifecyclePolicy(input)
	if err != nil {
		return err
	}

	// Save to database
	di := &db.Item{
		S3ObjectKey: u.item.record.S3.Object.Key,
		PolicyId:    *output.PolicyId,
		RequestId:   u.item.context.AwsRequestID,
		CreatedAt:   fmt.Sprintf("%s", u.item.record.EventTime),
		UpdatedAt:   fmt.Sprintf("%s", u.item.record.EventTime),
	}

	if err = u.dbconn.Create(di); err != nil {
		return err
	}

	return nil
}

// Update exsiting policy and related database record
func (u Upserter) UpdatePolicy() error {
	var input *dlm.UpdateLifecyclePolicyInput

	i, err := u.hydrate()
	if err != nil {
		return err
	}

	input, ok := i.(*dlm.UpdateLifecyclePolicyInput)
	if !ok {
		return errors.New("Failed to cast data into UpdateLifecyclePolicyInput")
	}

	// Update policy
	_, err = u.client.Dlm.UpdateLifecyclePolicy(input)
	if err != nil {
		return err
	}

	// Save to database
	di := &db.Item{
		S3ObjectKey: u.item.record.S3.Object.Key,
		PolicyId:    u.item.dbItem.PolicyId,
		RequestId:   u.item.context.AwsRequestID,
		CreatedAt:   u.item.dbItem.CreatedAt,
		UpdatedAt:   fmt.Sprintf("%s", u.item.record.EventTime),
	}

	if err = u.dbconn.Update(di); err != nil {
		return err
	}

	return nil
}

// Policy deleter
type Deleter struct {
	item   *eventItem
	client *AwsClients
	dbconn db.DB
}

// Strategy for delete
func (d Deleter) Execute() error {
	return d.DeletePolicy()
}

// Delete policy from DLM and related database record
func (d Deleter) DeletePolicy() error {
	if d.item.dbItem == nil {
		return fmt.Errorf("Failed to delete. No record has been found in database for policy %s", d.item.record.S3.Object.Key)
	}

	input := &dlm.DeleteLifecyclePolicyInput{
		PolicyId: aws.String(d.item.dbItem.PolicyId),
	}

	// Delete policy
	_, err := d.client.Dlm.DeleteLifecyclePolicy(input)
	if err != nil {
		return err
	}

	// Delete from database
	if err = d.dbconn.Delete(d.item.dbItem); err != nil {
		return err
	}

	return nil
}
