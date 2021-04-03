package model

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
	StartDate        *string   `json:"start_date"`
	EndDate          *string   `json:"end_date"`
	Title            string    `json:"title"`
	Company          string    `json:"company"`
	Responsibilities []*string `json:"responsibilities"`
}

func ExperienceTableName() *string {
	if val, ok := os.LookupEnv("SERVICE_TABLE_NAME"); ok {
		return &val
	}
	return aws.String("no tablename configured")
}

func ExperienceKey(lastname string, firstname string, company *string, title *string) map[string]*dynamodb.AttributeValue {
	sk := "experience"
	if company != nil {
		sk += fmt.Sprintf(" > company:%s", *company)
	}
	if title != nil {
		sk += fmt.Sprintf(" > title:%s", *title)
	}
	compKey := &CompositeKey{
		Pk: fmt.Sprintf("lName:%s > fName:%s", lastname, firstname),
		Sk: sk,
	}
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
		Key: ExperienceKey(input.LastName, input.FirstName, &input.Company, &input.Title),
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