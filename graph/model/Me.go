package model

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	clientmanager "github.com/entegral/aboutme/aws"
	e "github.com/entegral/aboutme/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	if _, ok := os.LookupEnv("SERVICE_TABLE_NAME"); !ok {
		logrus.Fatalln("SERVICE_TABLE_NAME env var must be set")
	}
}

type GetMeInput struct {
	FirstName       string        `json:"fName"`
	LastName        string        `json:"lName"`
}

type UpdateMeInput struct {
	FirstName       string        `json:"fName"`
	LastName        string        `json:"lName"`
	Title           *string       `json:"title,omitempty"`
	Location        *string       `json:"location,omitempty"`
	Interests       []*string     `json:"interests,omitempty"`
	AboutMe         []*string     `json:"about_me,omitempty"`
}

type Me struct {
	CompositeKey
	FirstName       string        `json:"fName"`
	LastName        string        `json:"lName"`
	Title           *string       `json:"title,omitempty"`
	Location        *string       `json:"location,omitempty"`
	Interests       []*string     `json:"interests,omitempty"`
	AboutMe         []*string     `json:"about_me,omitempty"`
	Experience      []*Experience `json:"experience,omitempty"`
	ExampleProjects []*Project    `json:"example_projects,omitempty"`
	Skills          *Skills       `json:"skills,omitempty"`
	Contact         *ContactInfo  `json:"contact,omitempty"`
}

// type helper functions
// these functions are used to ensure inputs and return types are aware of their storage configurations
func MeTableName() *string {
	if val, ok := os.LookupEnv("SERVICE_TABLE_NAME"); ok {
		return &val
	}
	return aws.String("no tablename configured")
}

func MeKey(lastname string, firstname string) map[string]*dynamodb.AttributeValue {
	compKey := &CompositeKey{
		Pk: fmt.Sprintf("lName:%s > fName:%s", lastname, firstname),
		Sk: fmt.Sprintf("general_info"),
	}
	key, err := dynamodbattribute.MarshalMap(compKey)
	if e.Warn("aboutme.Me.ReturnKey.MarshalMap", err) {
		return nil
	}
	return key
}

// Getter functions
// these functions receive input and return return types
func (input GetMeInput) GetMe() (*Me, error) {
	var m Me
	params := dynamodb.GetItemInput{
		TableName: MeTableName(),
		Key: MeKey(input.LastName, input.FirstName),
	}
	logrus.Warn("params", params)
	out, err := clientmanager.ClientMap["default"].Dynamo.GetItem(&params)
	if e.Warn("aboutme.Me.Get.GetItem", err) {
		return nil, err
	}
	logrus.Println(out.Item)
	if len(out.Item) == 0 {
		return nil, nil
	}
	err = dynamodbattribute.UnmarshalMap(out.Item, &m)
	if e.Warn("aboutme.Me.Get.UnmarshalMap", err) {
		return nil, err
	}
	return &m, err
}

func (m UpdateMeInput) Update() (*Me, error) {
	updateExpression := expression.Set(expression.Name("fName"), expression.Value(m.FirstName))
	updateExpression.Set(expression.Name("lName"), expression.Value(m.LastName))
	updateExpression.Set(expression.Name("title"), expression.Value(m.Title))
	updateExpression.Set(expression.Name("location"), expression.Value(m.Location))
	builder := expression.NewBuilder().WithUpdate(updateExpression)
	expr, err := builder.Build()
	if e.Warn("Me.Save.BuildExpression", err) {
		return nil, err
	}
	params := dynamodb.UpdateItemInput{
		Key: MeKey(m.LastName, m.FirstName),
		TableName: MeTableName(),
		UpdateExpression: expr.Update(),
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues: aws.String("ALL_NEW"),
	}
	logrus.Warn("params", params)
	var row Me
	out, err2 := clientmanager.ClientMap["default"].Dynamo.UpdateItem(&params)
	if e.Warn("Save.UpdateItem", err2) {
		return nil, err2
	}
	err3 := dynamodbattribute.UnmarshalMap(out.Attributes, &row)
	if e.Warn("error unmarshalling update response", err3) {
		return nil, err3
	}
	return &row, nil
}