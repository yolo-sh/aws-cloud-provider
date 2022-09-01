#!/bin/bash
# 
# Yolo base env AMI creation.
# 
set -euo pipefail

log () {
  echo -e "${1}" >&2
}

# Remove "debconf: unable to initialize frontend: Dialog" warnings
echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections

# -- Installing system dependencies

log "Installing system and Docker dependencies"

sudo apt-get --assume-yes --quiet --quiet update
sudo apt-get --assume-yes --quiet --quiet upgrade
sudo apt-get --assume-yes --quiet --quiet install jq curl wget git vim apt-transport-https ca-certificates gnupg lsb-release

# -- Installing Dokcer

log "Installing Docker"

if [[ ! -f "/usr/share/keyrings/docker-archive-keyring.gpg" ]]; then
  curl --fail --silent --show-error --location https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor --output /usr/share/keyrings/docker-archive-keyring.gpg
fi

if [[ ! -f "/etc/apt/sources.list.d/docker.list" ]]; then
	echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release --codename --short) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
fi

sudo apt-get --assume-yes --quiet --quiet update
sudo apt-get --assume-yes --quiet --quiet remove docker docker-engine docker.io containerd runc
sudo apt-get --assume-yes --quiet --quiet install docker-ce docker-ce-cli containerd.io

sudo docker pull yolosh/base-env:latest

sudo shred -u /etc/ssh/*_key /etc/ssh/*_key.pub
sudo rm --force /root/.ssh/authorized_keys
sudo rm --force /home/ubuntu/.ssh/authorized_keys
