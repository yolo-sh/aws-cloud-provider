#!/bin/bash
# 
# Yolo instance init.
# 
# This is the first script to run during the creation of the instance (via cloud-init).
# 
# See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html
# 
# In a nutshell, this script:
# - create and configure the user "yolo" (notably the SSH access)
# - configure and install the yolo agent
#
# The next steps are assured by the yolo agent via GRPC through SSH.
set -euo pipefail

log () {
  echo -e "${1}" >&2
}

log "\n\n"
log "---- Yolo instance init (start) ----"
log "\n\n"

# Remove "debconf: unable to initialize frontend: Dialog" warnings
echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections

# We use "jq" in our exit trap and "curl" to download the yolo agent
apt-get --assume-yes --quiet --quiet update
apt-get --assume-yes --quiet --quiet upgrade
apt-get --assume-yes --quiet --quiet install jq curl

constructExitJSONResponse () {
  JSON_RESPONSE=$(jq --null-input \
  --arg exitCode "${1}" \
  --arg sshHostKeys "${2}" \
  '{"exit_code": $exitCode, "ssh_host_keys": $sshHostKeys}')

  echo "${JSON_RESPONSE}"
}

YOLO_SSH_SERVER_HOST_KEY_FILE_PATH="/home/yolo/.ssh/yolo-ssh-server-host-key.pub"
YOLO_INIT_RESULTS_FILE_PATH="/tmp/yolo-init-results"

handleExit () {
  EXIT_CODE=$?

  rm --force "${YOLO_INIT_RESULTS_FILE_PATH}"

  log "\n\n"
  if [[ "${EXIT_CODE}" != 0 ]]; then
    constructExitJSONResponse "${EXIT_CODE}" "" >> "${YOLO_INIT_RESULTS_FILE_PATH}"
    log "---- Yolo instance init (failed) (exit code ${EXIT_CODE}) ----"
  else
    SSH_HOST_KEYS="$(cat "${YOLO_SSH_SERVER_HOST_KEY_FILE_PATH}")"
    constructExitJSONResponse "${EXIT_CODE}" "${SSH_HOST_KEYS}" >> "${YOLO_INIT_RESULTS_FILE_PATH}"
    
    log "---- Yolo instance init (success) ----"
  fi
  log "\n\n"

  exit "${EXIT_CODE}"
}

trap "handleExit" EXIT

# -- System configuration

# Lookup instance architecture for the yolo agent
INSTANCE_ARCH=""
case $(uname -m) in
  i386)       INSTANCE_ARCH="386" ;;
  i686)       INSTANCE_ARCH="386" ;;
  x86_64)     INSTANCE_ARCH="amd64" ;;
  arm)        dpkg --print-architecture | grep -q "arm64" && INSTANCE_ARCH="arm64" || INSTANCE_ARCH="armv6" ;;
  aarch64_be) INSTANCE_ARCH="arm64" ;;
  aarch64)    INSTANCE_ARCH="arm64" ;;
  armv8b)     INSTANCE_ARCH="arm64" ;;
  armv8l)     INSTANCE_ARCH="arm64" ;;
esac

# -- Create / Configure the user "yolo"

log "Creating user \"yolo\""

YOLO_CONFIG_DIR="/yolo-config"
YOLO_WORKSPACE_CONFIG_DIR="${YOLO_CONFIG_DIR}/workspace"

YOLO_USER_HOME_DIR="/home/yolo"

# Fixed gid because the yolo group is shared with container
groupadd --gid 10000 --force yolo
id -u yolo >/dev/null 2>&1 || useradd --gid yolo --home "${YOLO_USER_HOME_DIR}" --create-home --shell /bin/bash yolo

# Let the user "yolo" and the yolo agent
# run docker commands without "sudo".
# See https://docs.docker.com/engine/install/linux-postinstall/
groupadd --force docker
usermod --append --groups docker yolo

if [[ ! -f "/etc/sudoers.d/yolo" ]]; then
  echo "yolo ALL=(ALL) NOPASSWD:ALL" | tee /etc/sudoers.d/yolo > /dev/null
fi

mkdir --parents "${YOLO_CONFIG_DIR}"
mkdir --parents "${YOLO_WORKSPACE_CONFIG_DIR}"

chown --recursive yolo:yolo "${YOLO_CONFIG_DIR}"
chown --recursive yolo:yolo "${YOLO_USER_HOME_DIR}"

# Make sure that the user "yolo" in container
# (that share the same gid than the one in host)
# can write config folders
chmod --recursive 770 "${YOLO_CONFIG_DIR}"

log "Configuring home directory for user \"yolo\""

# We want the user "yolo" to be able to 
# connect through SSH via the generated SSH key.
# See below.
INSTANCE_SSH_PUBLIC_KEY="$(cat /home/ubuntu/.ssh/authorized_keys)"

# Run as "yolo"
sudo --set-home --login --user yolo -- env \
	INSTANCE_SSH_PUBLIC_KEY="${INSTANCE_SSH_PUBLIC_KEY}" \
bash << 'EOF'

mkdir --parents .ssh
chmod 700 .ssh

if [[ ! -f ".ssh/yolo-ssh-server-host-key" ]]; then
  ssh-keygen -t ed25519 -f .ssh/yolo-ssh-server-host-key -q -N ""
fi

chmod 644 .ssh/yolo-ssh-server-host-key.pub
chmod 600 .ssh/yolo-ssh-server-host-key

if [[ ! -f ".ssh/authorized_keys" ]]; then
  echo "${INSTANCE_SSH_PUBLIC_KEY}" >> .ssh/authorized_keys
fi

chmod 600 .ssh/authorized_keys

EOF

# -- Install the yolo agent
#
# /!\ the SSH server host key ("yolo-ssh-server-host-key")
#     needs to be generated. See above.
#
# /!\ the user "yolo" needs to be able to access 
#     the docker daemon. See above.

log "Installing the yolo agent"

YOLO_AGENT_VERSION="0.0.9"
YOLO_AGENT_TMP_ARCHIVE_PATH="/tmp/yolo-agent.tar.gz"
YOLO_AGENT_NAME="yolo-agent"
YOLO_AGENT_DIR="/usr/local/bin"
YOLO_AGENT_PATH="${YOLO_AGENT_DIR}/${YOLO_AGENT_NAME}"
YOLO_AGENT_SYSTEMD_SERVICE_NAME="yolo-agent.service"

if [[ ! -f "${YOLO_AGENT_PATH}" ]]; then
  #curl --fail --silent --show-error --location --header "Accept: application/octet-stream" https://api.github.com/repos/yolo-sh/agent/releases/assets/77939680 --output "${YOLO_AGENT_PATH}"
  rm --recursive --force "${YOLO_AGENT_TMP_ARCHIVE_PATH}"
  curl --fail --silent --show-error --location --header "Accept: application/octet-stream" "https://github.com/yolo-sh/agent/releases/download/v${YOLO_AGENT_VERSION}/agent_${YOLO_AGENT_VERSION}_linux_${INSTANCE_ARCH}.tar.gz" --output "${YOLO_AGENT_TMP_ARCHIVE_PATH}"
  tar --directory "${YOLO_AGENT_DIR}" --extract --file "${YOLO_AGENT_TMP_ARCHIVE_PATH}"
  rm --recursive --force "${YOLO_AGENT_TMP_ARCHIVE_PATH}"
fi

chmod +x "${YOLO_AGENT_PATH}"

if [[ ! -f "/etc/systemd/system/${YOLO_AGENT_SYSTEMD_SERVICE_NAME}" ]]; then
  tee /etc/systemd/system/"${YOLO_AGENT_SYSTEMD_SERVICE_NAME}" > /dev/null << EOF
  [Unit]
  Description=The agent used to establish connection with the Yolo CLI.

  [Service]
  Type=simple
  ExecStart=${YOLO_AGENT_PATH}
  WorkingDirectory=${YOLO_AGENT_DIR}
  Restart=always
  User=yolo
  Group=yolo

  [Install]
  WantedBy=multi-user.target
EOF
fi

systemctl enable "${YOLO_AGENT_SYSTEMD_SERVICE_NAME}"
systemctl start "${YOLO_AGENT_SYSTEMD_SERVICE_NAME}"
