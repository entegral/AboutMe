package aws

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type Clients struct {
	Dynamo		dynamodbiface.DynamoDBAPI
	SQS				sqsiface.SQSAPI
	S3				s3iface.S3API
}

var ClientMap = make(map[string]Clients)

func init() {
	CreateClients()
}
