package service_test

// import (
// 	"errors"
// 	"testing"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/golang/mock/gomock"
// 	"github.com/yolo-sh/aws-cloud-provider/mocks"
// 	"github.com/yolo-sh/aws-cloud-provider/service"
// 	"github.com/yolo-sh/aws-cloud-provider/userconfig"
// )

// func TestBuildWithResolvedUserConfig(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	resolvedUserConfig := userconfig.NewConfig("a", "b", "c")
// 	userConfigResolverMock := mocks.NewUserConfigResolver(mockCtrl)
// 	userConfigResolverMock.EXPECT().Resolve().Return(resolvedUserConfig, nil).AnyTimes()

// 	serviceConfigLoader := mocks.NewServiceConfigLoader(mockCtrl)
// 	serviceConfigLoader.EXPECT().Load(resolvedUserConfig).Return(aws.Config{}, nil)

// 	builder := service.NewBuilder(userConfigResolverMock, serviceConfigLoader)
// 	_, err := builder.Build()

// 	if err != nil {
// 		t.Fatalf("expected no error, got '%+v'", err)
// 	}

// 	// if sdkConfig.Region != resolvedUserConfig.Region {
// 	// 	t.Errorf("expected region to equal '%s', got '%s'", resolvedUserConfig.Region, sdkConfig.Region)
// 	// }

// 	// credsInSdk, err := sdkConfig.Credentials.Retrieve(context.TODO())

// 	// if err != nil {
// 	// 	t.Fatalf("expected no error, got '%+v'", err)
// 	// }

// 	// if credsInSdk.AccessKeyID != resolvedUserConfig.Credentials.AccessKeyID {
// 	// 	t.Errorf(
// 	// 		"expected access key id to equal '%s', got '%s'",
// 	// 		resolvedUserConfig.Credentials.AccessKeyID,
// 	// 		credsInSdk.AccessKeyID,
// 	// 	)
// 	// }

// 	// if credsInSdk.SecretAccessKey != resolvedUserConfig.Credentials.SecretAccessKey {
// 	// 	t.Errorf(
// 	// 		"expected secret access key to equal '%s', got '%s'",
// 	// 		resolvedUserConfig.Credentials.SecretAccessKey,
// 	// 		credsInSdk.SecretAccessKey,
// 	// 	)
// 	// }
// }

// func TestBuildWithUserConfigResolverError(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	userConfigResolveErr := userconfig.ErrMissingAccessKeyInEnv
// 	userConfigResolverMock := mocks.NewUserConfigResolver(mockCtrl)
// 	userConfigResolverMock.EXPECT().Resolve().Return(nil, userConfigResolveErr).AnyTimes()

// 	serviceConfigLoader := mocks.NewServiceConfigLoader(mockCtrl)

// 	builder := service.NewBuilder(userConfigResolverMock, serviceConfigLoader)
// 	_, err := builder.Build()

// 	if err == nil {
// 		t.Fatalf("expected error, got nothing")
// 	}

// 	if !errors.Is(err, userConfigResolveErr) {
// 		t.Fatalf("expected error to equal '%+v', got '%+v'", userConfigResolveErr, err)
// 	}
// }

// func TestBuildWithConfigLoaderError(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// 	resolvedUserConfig := userconfig.NewConfig("a", "b", "c")
// 	userConfigResolverMock := mocks.NewUserConfigResolver(mockCtrl)
// 	userConfigResolverMock.EXPECT().Resolve().Return(resolvedUserConfig, nil).AnyTimes()

// 	unknownError := errors.New("UnknownError")
// 	serviceConfigLoader := mocks.NewServiceConfigLoader(mockCtrl)
// 	serviceConfigLoader.EXPECT().Load(resolvedUserConfig).Return(aws.Config{}, unknownError)

// 	builder := service.NewBuilder(userConfigResolverMock, serviceConfigLoader)
// 	_, err := builder.Build()

// 	if err == nil {
// 		t.Fatalf("expected error, got nothing")
// 	}

// 	if !errors.Is(err, unknownError) {
// 		t.Fatalf("expected error to equal '%+v', got '%+v'", unknownError, err)
// 	}
// }
