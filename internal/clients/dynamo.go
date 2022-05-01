package clients

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	onceDynamo sync.Once
	dynamoConn *dynamodb.DynamoDB
)

func initializeDynamo() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamoConn = dynamodb.New(sess)
}

func DynamoDB() *dynamodb.DynamoDB {
	onceDynamo.Do(initializeDynamo)
	xray.AWS(dynamoConn.Client)
	return dynamoConn
}
