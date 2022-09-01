# Must match "PackerVariables" type

variable "aws_access_key_id" {
  type = string
}

variable "aws_secret_access_key" {
  type = string
}

variable "region" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "base_ami" {
  type = string
}
