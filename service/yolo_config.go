package service

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/yolo-sh/aws-cloud-provider/infrastructure"
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func (a *AWS) CreateYoloConfigStorage(
	stepper stepper.Stepper,
) error {

	dynamoDBClient := dynamodb.NewFromConfig(a.sdkConfig)

	stepper.StartTemporaryStep("Creating a DynamoDB table to store Yolo's data")

	err := infrastructure.CreateDynamoDBTableForYoloConfig(
		dynamoDBClient,
	)

	if err != nil && errors.Is(err, infrastructure.ErrYoloConfigTableAlreadyExists) {
		return nil
	}

	return err
}

func (a *AWS) LookupYoloConfig(
	stepper stepper.Stepper,
) (*entities.Config, error) {

	dynamoDBClient := dynamodb.NewFromConfig(a.sdkConfig)

	configJSON, err := infrastructure.LookupYoloConfigInDynamoDBTable(
		dynamoDBClient,
	)

	if err != nil {

		if errors.Is(err, infrastructure.ErrNoYoloConfigFound) {
			// No config table or no records.
			return nil, entities.ErrYoloNotInstalled
		}

		return nil, err
	}

	var yoloConfig *entities.Config
	err = json.Unmarshal([]byte(configJSON), &yoloConfig)

	if err != nil {
		return nil, err
	}

	return yoloConfig, nil
}

func (a *AWS) SaveYoloConfig(
	stepper stepper.Stepper,
	config *entities.Config,
) error {

	configJSON, err := json.Marshal(config)

	if err != nil {
		return err
	}

	dynamoDBClient := dynamodb.NewFromConfig(a.sdkConfig)

	return infrastructure.UpdateYoloConfigInDynamoDBTable(
		dynamoDBClient,
		config.ID,
		string(configJSON),
	)
}

func (a *AWS) RemoveYoloConfigStorage(
	stepper stepper.Stepper,
) error {

	dynamoDBClient := dynamodb.NewFromConfig(a.sdkConfig)

	stepper.StartTemporaryStep("Removing the DynamoDB table used to store Yolo's data")

	return infrastructure.RemoveDynamoDBTableForYoloConfig(
		dynamoDBClient,
	)
}
