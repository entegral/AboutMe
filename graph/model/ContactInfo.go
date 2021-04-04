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


type ContactInfo struct {
	Email           string        `json:"email"`
	LinkedIn        *string       `json:"linkedIn"`
	Github          *string       `json:"github"`
}

type ContactInfoInput struct {
	FirstName       string        `json:"fName"`
	LastName        string        `json:"lName"`
	Email           *string        `json:"email"`
	LinkedIn        *string       `json:"linkedIn"`
	Github          *string       `json:"github"`
}

func ContactInfoTableName() *string {
	if val, ok := os.LookupEnv("SERVICE_TABLE_NAME"); ok {
		return &val
	}
	return aws.String("no tablename configured")
}

func ContactInfoKey(lastname string, firstname string) map[string]*dynamodb.AttributeValue {
	compKey := &CompositeKey{
		Pk: fmt.Sprintf("lName:%s > fName:%s", lastname, firstname),
		Sk: fmt.Sprintf("contact_info"),
	}
	key, err := dynamodbattribute.MarshalMap(compKey)
	if e.Warn("ContactInfo.ContactInfoKey.MarshalMap", err) {
		return nil
	}
	return key
}


func (input ContactInfoInput) Get() (*ContactInfo, error) {
	var c ContactInfo
	params := dynamodb.GetItemInput{
		TableName: ContactInfoTableName(),
		Key: ContactInfoKey(input.LastName, input.FirstName),
	}
	logrus.Debug("params", params)
	out, err := clientmanager.ClientMap["default"].Dynamo.GetItem(&params)
	if e.Warn("error in ContactInfo.Get.GetItem", err) {
		return nil, err
	}
	logrus.Debug(out.Item)
	if len(out.Item) == 0 {
		return nil, nil
	}
	err = dynamodbattribute.UnmarshalMap(out.Item, &c)
	if e.Warn("error in ContactInfo.Get.UnmarshalMap", err) {
		return nil, err
	}
	return &c, err
}

func (input ContactInfoInput) Update() (*ContactInfo, error) {
	updateExpression := expression.Set(expression.Name("fName"), expression.Value(input.FirstName))
	updateExpression.Set(expression.Name("lName"), expression.Value(input.LastName))
	if input.Email != nil && *input.Email != "" {
		updateExpression.Set(expression.Name("email"), expression.Value(input.Email))
	}
	if input.LinkedIn != nil && *input.LinkedIn != "" {
		updateExpression.Set(expression.Name("linkedIn"), expression.Value(input.LinkedIn))
	}
	if input.Github != nil && *input.Github != "" {
		updateExpression.Set(expression.Name("github"), expression.Value(input.Github))
	}
	builder := expression.NewBuilder().WithUpdate(updateExpression)
	expr, err := builder.Build()
	if e.Warn("ContactInfo.Update.BuildExpression", err) {
		return nil, err
	}
	params := dynamodb.UpdateItemInput{
		Key: ContactInfoKey(input.LastName, input.FirstName),
		TableName: ContactInfoTableName(),
		UpdateExpression: expr.Update(),
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues: aws.String("ALL_NEW"),
	}
	logrus.Debug("params", params)
	var row ContactInfo
	out, err2 := clientmanager.ClientMap["default"].Dynamo.UpdateItem(&params)
	if e.Warn("error updating ContactInfo record", err2) {
		return nil, err2
	}
	err3 := dynamodbattribute.UnmarshalMap(out.Attributes, &row)
	if e.Warn("error unmarshalling update response for contact info record", err3) {
		return nil, err3
	}
	return &row, nil

}