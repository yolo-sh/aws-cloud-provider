package service

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/yolo-sh/aws-cloud-provider/infrastructure"
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func (a *AWS) OpenPort(
	stepper stepper.Stepper,
	config *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
	portToOpen string,
) error {

	var envInfra *EnvInfrastructure
	err := json.Unmarshal([]byte(env.InfrastructureJSON), &envInfra)

	if err != nil {
		return err
	}

	ec2Client := ec2.NewFromConfig(a.sdkConfig)

	err = infrastructure.OpenInstancePort(
		ec2Client,
		envInfra.SecurityGroup.ID,
		portToOpen,
	)

	if err != nil {
		return err
	}

	env.OpenedPorts[portToOpen] = true

	return nil
}
