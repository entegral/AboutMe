package model

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	clientmanager "github.com/entegral/aboutme/aws"
	e "github.com/entegral/aboutme/errors"
	"github.com/sirupsen/logrus"
)
type Experience struct {
	FirstName       string        `json:"fName"`
	LastName        string        `json:"lName"`
	StartDate        *string   `json:"start_date"`
	EndDate          *string   `json:"end_date"`
	Title            *string   `json:"title"`
	Company          *string   `json:"company"`
	Responsibilities []*string `json:"responsibilities"`
}

type ExperienceInput struct {
	FirstName       string        `json:"fName"`
	LastName        string        `json:"lName"`
	StartDate       string   `json:"start_date"`
	EndDate          *string   `json:"end_date"`
	Title            *string    `json:"title"`
	Company          *string    `json:"company"`
	Responsibilities []*string `json:"responsibilities"`
}

func ExperienceTableName() *string {
	if val, ok := os.LookupEnv("SERVICE_TABLE_NAME"); ok {
		return &val
	}
	return aws.String("no tablename configured")
}

func ExperienceCompKey(lastname string, firstname string, startDate string) CompositeKey {
	return CompositeKey{
		Pk: fmt.Sprintf("lName:%s > fName:%s", lastname, firstname),
		Sk: fmt.Sprintf("experience > startDate:%v", startDate),
	}
}
func ExperienceKey(lastname string, firstname string, startDate string) map[string]*dynamodb.AttributeValue {
	compKey := ExperienceCompKey(lastname, firstname, startDate)
	key, err := dynamodbattribute.MarshalMap(compKey)
	if e.Warn("error occurred during marshalling of experience key", err) {
		return nil
	}
	return key
}

func (input ExperienceInput) Get() (*Experience, error) {
	var ex Experience
	params := dynamodb.GetItemInput{
		TableName: MeTableName(),
		Key: ExperienceKey(input.LastName, input.FirstName, input.StartDate),
	}
	logrus.Warn("params", params)
	out, err := clientmanager.ClientMap["default"].Dynamo.GetItem(&params)
	if e.Warn("error in Experience.Get.GetItem", err) {
		return nil, err
	}
	logrus.Println(out.Item)
	if len(out.Item) == 0 {
		return nil, nil
	}
	err = dynamodbattribute.UnmarshalMap(out.Item, &ex)
	if e.Warn("error in Experience.Get.UnmarshalMap", err) {
		return nil, err
	}
	return &ex, err
}

func convertDateFormat(date string) (*time.Time, error) {
	t, err := time.Parse(layoutUS, date)
	if err != nil {
		return nil, errors.New("Start/End dates must be provided in following format: 'Month Day, Year'")
	}
	return &t, nil
}

func (input ExperienceInput) Query() ([]*Experience, error) {
	var experiences []*Experience

	compKey := ExperienceCompKey(input.LastName, input.FirstName, input.StartDate)
	params := dynamodb.QueryInput{
		TableName: MeTableName(),
		KeyConditionExpression: aws.String("pk = :pk and begins_with(sk, :sk)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: &compKey.Pk},
			":sk": {S: &compKey.Sk},
		},
		ScanIndexForward: aws.Bool(false),
	}
	logrus.Warn("params", params)
	out, err := clientmanager.ClientMap["default"].Dynamo.Query(&params)
	if e.Warn("error in Experience.Get.Query", err) {
		return nil, err
	}
	logrus.Println(out.Items)
	if len(out.Items) == 0 {
		return nil, nil
	}
	for i, item := range out.Items {
		var ex Experience
		err = dynamodbattribute.UnmarshalMap(item, &ex)
		if e.Warn("error in Experience.Query.UnmarshalMap loop", err) {
			logrus.Println("error item:", item)
			logrus.Println("error iterant:", i)
		}
		experiences = append(experiences, &ex)
	}
	logrus.Warn("Found experiences", experiences)
	return experiences, err
}
const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

func (input ExperienceInput) Update() (*Experience, error) {
	updateExpression := expression.Set(expression.Name("fName"), expression.Value(input.FirstName))
	updateExpression.Set(expression.Name("lName"), expression.Value(input.LastName))
	updateExpression.Set(expression.Name("start_date"), expression.Value(input.StartDate))
	if input.Title != nil && *input.Title != "" {
		updateExpression.Set(expression.Name("title"), expression.Value(input.Title))
	}
	if input.Company != nil && *input.Company != "" {
		updateExpression.Set(expression.Name("company"), expression.Value(input.Company))
	}
	if len(input.Responsibilities) != 0 {
		updateExpression.Set(expression.Name("responsibilities"), expression.Value(input.Responsibilities))
	}
	if input.EndDate != nil &&  *input.EndDate != "" {
		_, endDateErr := convertDateFormat(*	input.EndDate)
		if e.Warn("error parsing end date", endDateErr) {
			return nil, endDateErr
		}
		updateExpression.Set(expression.Name("end_date"), expression.Value(input.EndDate))
	}
	builder := expression.NewBuilder().WithUpdate(updateExpression)
	expr, err := builder.Build()
	if e.Warn("Me.Update.BuildExpression", err) {
		return nil, err
	}
	t, timeErr := convertDateFormat(input.StartDate)
	if e.Warn("bad start time input provided", timeErr) {
		return nil, timeErr
	}
	params := dynamodb.UpdateItemInput{
		Key: ExperienceKey(input.LastName, input.FirstName, t.Format(layoutISO)),
		TableName: MeTableName(),
		UpdateExpression: expr.Update(),
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues: aws.String("ALL_NEW"),
	}
	logrus.Warn("params", params)
	var row Experience
	out, err2 := clientmanager.ClientMap["default"].Dynamo.UpdateItem(&params)
	if e.Warn("error updating Experience record", err2) {
		return nil, err2
	}
	err3 := dynamodbattribute.UnmarshalMap(out.Attributes, &row)
	if e.Warn("error unmarshalling update response for experience record", err3) {
		return nil, err3
	}
	return &row, nil
}