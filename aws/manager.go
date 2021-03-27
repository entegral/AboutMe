package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

func CreateClients() Clients {
	logrus.Debugln("Creating default account clients")
	sess := session.Must(session.NewSession())
	var c Clients
	c.Dynamo = dynamodb.New(sess)
	c.SQS = sqs.New(sess)
	c.S3 = s3.New(sess)
	ClientMap["default"] = c
	logrus.Debugln("Default account clients added:")
	return c
}

func CreateCustomClientConnection(region string, profileName string) Clients {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile:           profileName,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		logrus.Fatalf("session init failed with AWS_PROFILE: %+v, message: %s", profileName, err.Error())
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		logrus.Fatalf("bad session credentials for AWS_PROFILE: %+v, message: %s", profileName, err.Error())
	}
	cfg := &aws.Config{Credentials: sess.Config.Credentials}
	var c Clients
	c.Dynamo = dynamodb.New(sess, cfg)
	c.SQS = sqs.New(sess, cfg)
	c.S3 = s3.New(sess, cfg)
	return c
}

func AddConnection(accountName string, client Clients) {
	if _, ok := ClientMap[accountName]; !ok {
		ClientMap[accountName] = client
		logrus.Debugln("client added for account: %s", accountName)
	}
}
