package config

import (
	"regexp"

	"github.com/yolo-sh/aws-cloud-provider/userconfig"

	awsRegions "github.com/jsonmaur/aws-regions/v2"
)

const (
	AWSAccessKeyIDPattern     = "^[A-Z0-9]{20}$"
	AWSSecretAccessKeyPattern = "^[A-Za-z0-9/+=]{40}$"
)

type UserConfigValidator struct{}

func NewUserConfigValidator() UserConfigValidator {
	return UserConfigValidator{}
}

func (u UserConfigValidator) Validate(userConfig *userconfig.Config) error {
	region := userConfig.Region

	if err := u.validateRegion(region); err != nil {
		return err
	}

	creds := userConfig.Credentials
	accessKeyID := creds.AccessKeyID
	secretAccessKey := creds.SecretAccessKey

	if err := u.validateAccessKeyID(accessKeyID); err != nil {
		return err
	}

	if err := u.validateSecretAccessKey(secretAccessKey); err != nil {
		return err
	}

	return nil
}

func (UserConfigValidator) validateRegion(region string) error {
	_, err := awsRegions.LookupByCode(region)

	if err != nil {
		return ErrInvalidRegion{
			Region: region,
		}
	}

	return nil
}

func (UserConfigValidator) validateAccessKeyID(accessKeyID string) error {
	match, err := regexp.MatchString(AWSAccessKeyIDPattern, accessKeyID)

	if err != nil {
		return err
	}

	if !match {
		return ErrInvalidAccessKeyID{
			AccessKeyID: accessKeyID,
		}
	}

	return nil
}

func (UserConfigValidator) validateSecretAccessKey(secretAccessKey string) error {
	match, err := regexp.MatchString(AWSSecretAccessKeyPattern, secretAccessKey)

	if err != nil {
		return err
	}

	if !match {
		return ErrInvalidSecretAccessKey{
			SecretAccessKey: secretAccessKey,
		}
	}

	return nil
}
