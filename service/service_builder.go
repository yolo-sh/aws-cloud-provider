package service

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/yolo-sh/aws-cloud-provider/userconfig"
	"github.com/yolo-sh/yolo/entities"
)

//go:generate mockgen -destination ../mocks/user_config_resolver.go -package mocks -mock_names UserConfigResolver=UserConfigResolver github.com/yolo-sh/aws-cloud-provider/service UserConfigResolver
type UserConfigResolver interface {
	Resolve() (*userconfig.Config, error)
}

//go:generate mockgen -destination ../mocks/user_config_loader.go -package mocks -mock_names ConfigLoader=UserConfigLoader github.com/yolo-sh/aws-cloud-provider/service UserConfigLoader
type UserConfigLoader interface {
	Load(userConfig *userconfig.Config) (aws.Config, error)
}

type UserConfigValidator interface {
	Validate(userConfig *userconfig.Config) error
}

type Builder struct {
	userConfigResolver  UserConfigResolver
	userConfigValidator UserConfigValidator
	userConfigLoader    UserConfigLoader
}

func NewBuilder(
	userConfigResolver UserConfigResolver,
	userConfigValidator UserConfigValidator,
	userConfigLoader UserConfigLoader,
) Builder {

	return Builder{
		userConfigResolver:  userConfigResolver,
		userConfigValidator: userConfigValidator,
		userConfigLoader:    userConfigLoader,
	}
}

func (b Builder) Build() (entities.CloudService, error) {
	userConfig, err := b.userConfigResolver.Resolve()

	if err != nil {
		return nil, err
	}

	if err := b.userConfigValidator.Validate(userConfig); err != nil {
		return nil, err
	}

	AWSSDKConfig, err := b.userConfigLoader.Load(userConfig)

	if err != nil {
		return nil, err
	}

	AWSService := NewAWS(AWSSDKConfig)

	return AWSService, nil
}
