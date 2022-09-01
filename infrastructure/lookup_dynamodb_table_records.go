package infrastructure

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	ErrNoYoloConfigFound       = errors.New("ErrNoYoloConfigFound")
	ErrMultipleYoloConfigFound = errors.New("ErrMultipleYoloConfigFound")
)

type DynamoDBYoloConfigTableRecord struct {
	ID         string
	ConfigJSON string
}

func LookupYoloConfigInDynamoDBTable(
	dynamoDBClient *dynamodb.Client,
) (returnedConfigJSON string, returnedError error) {

	scanResp, err := dynamoDBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(DynamoDBYoloConfigTableName),
	})

	if err != nil {
		var resourceNotFoundErr *types.ResourceNotFoundException

		if errors.As(err, &resourceNotFoundErr) { // Table not found
			returnedError = ErrNoYoloConfigFound
			return
		}

		returnedError = err
		return
	}

	if scanResp.Count == 0 { // Empty table
		returnedError = ErrNoYoloConfigFound
		return
	}

	if scanResp.Count > 1 { // Multiple rows
		returnedError = ErrMultipleYoloConfigFound
		return
	}

	var records []DynamoDBYoloConfigTableRecord
	err = attributevalue.UnmarshalListOfMaps(scanResp.Items, &records)

	if err != nil {
		returnedError = err
		return
	}

	returnedConfigJSON = records[0].ConfigJSON
	return
}
