package pkg

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var sess, err = session.NewSession(&aws.Config{})

// DynamoDB table names passed via environment
var ddbTableProducts string
var ddbTableCategories string

// Allow DDB endpoint to be overridden to support amazon/dynamodb-local
var ddbEndpointOverride string
var runningLocal bool

var dynamoClient *dynamodb.DynamoDB

// Initialize clients
func InitProduct() {
	if len(ddbEndpointOverride) > 0 {
		runningLocal = true
		log.Println("Creating DDB client with endpoint override: ", ddbEndpointOverride)
		creds := credentials.NewStaticCredentials("does", "not", "matter")
		awsConfig := &aws.Config{
			Credentials: creds,
			Region:      aws.String("ap-northeast-2"),
			Endpoint:    aws.String(ddbEndpointOverride),
		}
		dynamoClient = dynamodb.New(sess, awsConfig)
	} else {
		runningLocal = false
		dynamoClient = dynamodb.New(sess)
	}
}
