package model

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	clientmanager "github.com/entegral/aboutme/aws"
	"github.com/entegral/aboutme/errors"
	"github.com/sirupsen/logrus"
)

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

type CompositeKey struct {
	Pk	string		`json:"pk"`
	Sk	string		`json:"sk"`
}

func (m Me) TableName() *string {
	if val, ok := os.LookupEnv("SERVICE_TABLE_NAME"); ok {
		return &val
	}
	return aws.String("no tablename configured")
}

func returnKey(lastName string, firstName string) map[string]*dynamodb.AttributeValue {
	compKey := &CompositeKey{
		Pk: fmt.Sprintf("lName:%s > fName:%s", lastName, firstName),
		Sk: fmt.Sprintf("general_info"),
	}
	key, err := dynamodbattribute.MarshalMap(compKey)
	if errors.Warn("aboutme.Me.ReturnKey.MarshalMap", err) {
		return nil
	}
	return key
}

func (c CompositeKey) GetMe(lastName string, firstName string) (*Me, error) {
	var m Me
	params := dynamodb.GetItemInput{
		TableName: m.TableName(),
		Key: returnKey(lastName, firstName),
	}
	logrus.Warn("params", params)
	out, err := clientmanager.ClientMap["default"].Dynamo.GetItem(&params)
	if errors.Warn("aboutme.Me.Get.GetItem", err) {
		return nil, err
	}
	logrus.Println(out.Item)
	if len(out.Item) == 0 {
		return nil, nil
	}
	err = dynamodbattribute.UnmarshalMap(out.Item, &m)
	if errors.Warn("aboutme.Me.Get.UnmarshalMap", err) {
		return nil, err
	}
	return &m, err
}

func (m Me) Save() (*Me, error) {
	updateExpression := expression.Set(expression.Name("fName"), expression.Value(m.FirstName))
	updateExpression.Set(expression.Name("lName"), expression.Value(m.LastName))
	updateExpression.Set(expression.Name("title"), expression.Value(m.Title))
	updateExpression.Set(expression.Name("location"), expression.Value(m.Location))
	builder := expression.NewBuilder().WithUpdate(updateExpression)
	expr, err := builder.Build()
	if errors.Warn("Me.Save.BuildExpression", err) {
		return nil, err
	}
	params := dynamodb.UpdateItemInput{
		Key: returnKey(m.LastName, m.FirstName),
		TableName: m.TableName(),
		UpdateExpression: expr.Update(),
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues: aws.String("ALL_NEW"),
	}
	logrus.Warn("params", params)
	_, err2 := clientmanager.ClientMap["default"].Dynamo.UpdateItem(&params)
	if errors.Warn("Save.UpdateItem", err2) {
		return nil, err2
	}
	return &m, nil
}