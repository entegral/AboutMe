package model

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	clientmanager "github.com/entegral/aboutme/aws"
	"github.com/entegral/aboutme/errors"
	"github.com/sirupsen/logrus"
)

type Me struct {
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

func (c CompositeKey) GetMe(lastName string, firstName string) (*Me, error) {
	var m Me
	compKey := &CompositeKey{
		Pk: fmt.Sprintf("lName:%s > fName:%s", lastName, firstName),
		Sk: fmt.Sprintf("general_info"),
	}
	key, err := dynamodbattribute.MarshalMap(compKey)
	if errors.Warn("aboutme.Me.Get.MarshalMap", err) {
		return nil, err
	}
	params := dynamodb.GetItemInput{
		TableName: m.TableName(),
		Key: key,
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