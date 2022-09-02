<p align="center">
  <img src="https://user-images.githubusercontent.com/1233275/172015952-8f40a3e0-b07b-4c57-8526-bc0870aad76d.jpeg" alt="yolo" width="180" height="180" />
</p>

<p align="center">
    <h1 align="center">AWS Cloud Provider</h1>
    <p align="center">This repository contains the source code that implements the AWS cloud provider for the <a href="https://github.com/yolo-sh/cli">Yolo CLI</a>.</p>
</p>

```bash
yolo aws --profile production --region eu-west-3 init yolo-sh/api
```
<p align="center">
  <img width="759" src="https://user-images.githubusercontent.com/1233275/187968190-0ce41e41-4612-486e-bc3d-64f03fc25ac5.png" alt="Example of use of the Yolo CLI" />
</p>

## Table of contents
- [Usage](#usage)
    - [Authentication](#authentication)
        - [--profile](#--profile)
        - [--region](#--region)
    - [Permissions](#permissions)
    - [Authorized instance types](#authorized-instance-types)
- [Infrastructure components](#infrastructure-components)
    - [Init](#init)
    - [Edit](#edit)
    - [Open port](#open-port)
    - [Close port](#close-port)
    - [Remove](#remove)
    - [Uninstall](#uninstall)
- [Infrastructure costs](#infrastructure-costs)
- [License](#license)

## Usage

```console
To begin, create your first environment using the command:

  yolo aws init <repository>

Once initialized, you may want to connect to it using the command: 

  yolo aws edit <repository>

If you don't plan to use this environment again, you could remove it using the command:
	
  yolo aws remove <repository>

<repository> may be relative to your personal GitHub account (eg: cli) or fully qualified (eg: my-organization/api).

Usage:
  yolo aws [command]

Examples:
  yolo aws init yolo-sh/api --instance-type m4.large 
  yolo aws edit yolo-sh/api
  yolo aws remove yolo-sh/api

Available Commands:
  close-port  Close a port in an environment
  edit        Connect your editor to an environment
  init        Initialize a new environment
  open-port   Open a port in an environment
  remove      Remove an environment
  uninstall   Uninstall Yolo from your AWS account

Flags:
  -h, --help             help for aws
      --profile string   the configuration profile to use to access your AWS account
      --region string    the region to use to access your AWS account

Use "yolo aws [command] --help" for more information about a command.
```

### Authentication

In order to access your AWS account, the Yolo CLI will first look for credentials in the following environment variables:

- `AWS_ACCESS_KEY_ID`

- `AWS_SECRET_ACCESS_KEY`

- `AWS_REGION`

If not found, the configuration files created by the AWS CLI (via `aws configure`) will be used.

#### --profile

If you have configured the AWS CLI with multiple configuration profiles, you could tell Yolo which one to use via the `--profile` flag:

```shell
yolo aws --profile production init yolo-sh/api
```

**By default, Yolo will use the profile named `default`.**

#### --region

If you want to overwrite the region resolved by the Yolo CLI, you could use the `--region` flag:

```shell
yolo aws --region eu-west-3 init yolo-sh/api
```

```shell
yolo aws --profile production --region eu-west-3 init yolo-sh/api
```

### Permissions

Your credentials must have certain permissions attached to be used with Yolo. See the next sections to learn more about the actions that will be done on your behalf.

### Authorized instance types

To be used with Yolo, the chosen instance must be **an on-demand linux instance (with EBS support) running on an amd64 or arm64 architecture**.

#### Examples

```shell
t2.medium, m6g.large, a1.xlarge, c5.12xlarge...
```

## Infrastructure components

![infra](https://user-images.githubusercontent.com/1233275/187925670-e06790b5-0084-4d91-a18e-160c771b4f4a.png)

The schema above describe all the components that may be created in your AWS account. The next sections will describe their lifetime according to your use of the Yolo CLI.

### Init

```bash
yolo aws init yolo-sh/api --instance-type t2.medium
```

#### The first time Yolo is used in a region

A DynamoDB table named `yolo-config-dynamodb-table` will be created. This table will be used to store the state of the Yolo's infrastructure.

Once created, all the following components will also be created:

- A `VPC` named `yolo-vpc` with an IPv4 CIDR block equals to `10.0.0.0/16` to isolate your infrastructure.

- A `public subnet` named `yolo-public-subnet` with an IPv4 CIDR block equals to `10.0.0.0/24` that will contain your environments' instances.

- An `internet gateway` named `yolo-internet-gateway` to let your environments' instances communicate with internet.

- A `route table` named `yolo-route-table` that will allow egress traffic from your your environments' instances to the internet (via the internet gateway).

#### On each init

Each time the `init` command is run for a new environment, the following components will be created:

- A `security group` named `yolo-${ENV_NAME}-security-group` to let your environment accepts `SSH` connections on port `2200`.

- An `SSH key pair` named `yolo-${ENV_NAME}-key-pair` to let you access your environment via `SSH`.

- A `network interface` named `yolo-${ENV_NAME}-network-interface` to enable network connectivity in your environment.

- An `Elastic IP` named `yolo-${ENV_NAME}-elastic-ip` to let you access your environment via a fixed public IP.

- An `EC2 instance` named `yolo-${ENV_NAME}-instance` with a type equals to the one passed via the `--instance-type` flag or `t2.medium` by default.
    
- An `EBS volume` attached to the instance (default to `16GB`) to provide long-term storage to your environment.

### Edit

```bash
yolo aws edit yolo-sh/api
```

When running the `edit` command, nothing will be done to your infrastructure.

### Open port

```bash
yolo aws open-port yolo-sh/api 8080
```

When running the `open-port` command, an `ingress` rule will be added to the `security group` of the environment. 

This rule will allow all `TCP` trafic from `any IP address` to the specified port.

### Close port

```bash
yolo aws close-port yolo-sh/api 8080
```

When running the `close-port` command, the `ingress` rule added by the `open-port` command will be removed.

### Remove

```bash
yolo aws remove yolo-sh/api
```

When running the `remove` command, all the components associated with the environment will be removed.

In other words:

- The `EC2 instance`.

- The `Elastic IP`.

- The `network interface`.

- The `EBS volume`.

- The `SSH key pair`.

- The `security group`.

### Uninstall

```bash
yolo aws uninstall
```

When running the `uninstall` command, all the shared components will be removed. 

In other words:

- The `route table`.

- The `internet gateway`.

- The `public subnet`.

- The `VPC`.

- The `DynamoDB table`.

## Infrastructure costs

The costs of running an environment on AWS are essentially equal to the costs of the `EC2` instance and the `EBS` volume:

- For the `EC2` instance, the price depends on the instance type chosen.

- For the `EBS` volume, Yolo uses the `General Purpose SSD (gp2)` type that will cost you ~$0.10 per GB-month.

All other components are free (or mostly free) given their limited usage.

## License

Yolo is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
