package model

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	clientmanager "github.com/entegral/aboutme/aws"
	"github.com/entegral/aboutme/errors"
)

type Me struct {
	FirstName       string        `json:"first_name"`
	LastName        string        `json:"last_name"`
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

func (c CompositeKey) GetMe(lastName string, firstName string) (Me, error) {
	var m Me
	compKey := CompositeKey{
		Pk: fmt.Sprintf("lName > %s > fName > %s", lastName, firstName),
		Sk: fmt.Sprintf("general_info"),
	}
	key, err := dynamodbattribute.MarshalMap(compKey)
	if errors.Warn("aboutme.Me.Get.MarshalMap", err) {
		return m, err
	}
	params := dynamodb.GetItemInput{
		TableName: m.TableName(),
		Key: key,
	}
	out, err := clientmanager.ClientMap["default"].Dynamo.GetItem(&params)
	if errors.Warn("aboutme.Me.Get.GetItem", err) {
		return m, err
	}
	err = dynamodbattribute.UnmarshalMap(out.Item, m)
	if errors.Warn("aboutme.Me.Get.UnmarshalMap", err) {
		return m, err
	}
	return m, err
}