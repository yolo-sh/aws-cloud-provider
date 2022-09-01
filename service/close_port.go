package service

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/yolo-sh/aws-cloud-provider/infrastructure"
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func (a *AWS) ClosePort(
	stepper stepper.Stepper,
	config *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
	portToClose string,
) error {

	var envInfra *EnvInfrastructure
	err := json.Unmarshal([]byte(env.InfrastructureJSON), &envInfra)

	if err != nil {
		return err
	}

	ec2Client := ec2.NewFromConfig(a.sdkConfig)

	err = infrastructure.CloseInstancePort(
		ec2Client,
		envInfra.SecurityGroup.ID,
		portToClose,
	)

	if err != nil {
		return err
	}

	delete(env.OpenedPorts, portToClose)

	return nil
}
