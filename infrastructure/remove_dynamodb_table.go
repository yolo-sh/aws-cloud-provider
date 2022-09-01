package infrastructure

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func RemoveDynamoDBTableForYoloConfig(
	dynamoDBClient *dynamodb.Client,
) error {

	_, err := dynamoDBClient.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(DynamoDBYoloConfigTableName),
	})

	if err != nil {
		return err
	}

	waiter := dynamodb.NewTableNotExistsWaiter(dynamoDBClient)
	maxWaitTime := 5 * time.Minute

	return waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(DynamoDBYoloConfigTableName),
	}, maxWaitTime)
}
