package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const DynamoDBYoloConfigTableName = "yolo-configuration-dynamodb-table"

var (
	ErrYoloConfigTableAlreadyExists = errors.New("ErrYoloConfigTableAlreadyExists")
)

func CreateDynamoDBTableForYoloConfig(
	dynamoDBClient *dynamodb.Client,
) error {

	_, err := dynamoDBClient.CreateTable(
		context.TODO(),

		&dynamodb.CreateTableInput{
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("ID"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},

			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("ID"),
					KeyType:       types.KeyTypeHash,
				},
			},

			BillingMode: types.BillingModePayPerRequest,

			TableName: aws.String(DynamoDBYoloConfigTableName),
		},
	)

	if err != nil {
		var tableExistsError *types.ResourceInUseException

		if errors.As(err, &tableExistsError) {
			return ErrYoloConfigTableAlreadyExists
		}

		return err
	}

	existsWaiter := dynamodb.NewTableExistsWaiter(dynamoDBClient)
	maxWaitTime := 5 * time.Minute

	return existsWaiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(DynamoDBYoloConfigTableName),
	}, maxWaitTime)
}
