locals {
  resources_name_prefix = "yolo_base_env"
  timestamp             = regex_replace(timestamp(), "[- TZ:]", "")
}

packer {
  required_plugins {
    amazon = {
      version = ">= 1.0.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "base_image" {
  access_key = var.aws_access_key_id
  secret_key = var.aws_secret_access_key
  region     = var.region

  ami_name        = "${local.resources_name_prefix}_base_ami_${local.timestamp}"
  ami_description = "The AMI used as base by your newly created VMs"

  instance_type = var.instance_type

  source_ami_filter {
    filters = {
      name                = var.base_ami
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }

  ssh_username            = "ubuntu"
  temporary_key_pair_type = "ed25519"

  launch_block_device_mappings {
    device_name           = "/dev/sda1"
    volume_size           = 16
    volume_type           = "gp2"
    delete_on_termination = true
  }

  tags = {
    Name = "${local.resources_name_prefix}_base_ami_${local.timestamp}"
  }

  run_tags = {
    Name = "${local.resources_name_prefix}_base_ami_builder"
  }
}

build {
  sources = ["source.amazon-ebs.base_image"]

  # Wait for instance initialization to prevent apt errors
  # See https://www.packer.io/docs/debugging#issues-installing-ubuntu-packages
  provisioner "shell" {
    inline = ["/usr/bin/cloud-init status --wait"]
  }

  provisioner "shell" {
    script = "./bootstrap.sh"
  }

  // Create a "packer-manifest.json" file with build infos
  // Used to access AMI ID
  post-processor "manifest" {}
}
