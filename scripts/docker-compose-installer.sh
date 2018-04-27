#!/usr/bin/env bash
VERSION="v1.0.0"
DOCKER_COMPOSE_CMD="docker-compose"
DOCKER_COMPOSE_VERSION="1.21.0"
DOWNLOAD_URL_DOCKER_COMPOSE="https://raw.githubusercontent.com/mmichaelb/gosharexserver/master/deployments/docker-compose.yml"
DOWNLOAD_URL_DOCKER_COMPOSE_CONFIG="https://raw.githubusercontent.com/mmichaelb/gosharexserver/master/configs/docker-compose-config.toml"
echo "gosharexserver docker-compose installer ${VERSION} (https://github.com/mmichaelb/gosharexserver)"
echo
echo "Checking for docker-compose installation..."
DOCKER_COMPOSE_VERSION=$(${DOCKER_COMPOSE_CMD} version --short)
if [[ ${DOCKER_COMPOSE_VERSION} =~ [0-9]+.[0-9]+.[0-9]+ ]]; then
    echo "Found existing version of docker compose (v${DOCKER_COMPOSE_VERSION})."
else
    read -r -p "No Docker compose version found. Should a new version be installed? [y/n] " REPLY
    case $REPLY in
    [yY])
        echo "Installing Docker compose..."
        sudo curl -L https://github.com/docker/compose/releases/download/1.21.0/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
        echo "Done with installation of Docker compose (v${DOCKER_COMPOSE_VERSION})."
        ;;
    *)
        echo "Docker compose is needed for the installation. If you use a different command name, adjust the DOCKER_COMPOSE_CMD value."
        exit 1
        ;;
    esac
fi
echo "Downloading data from git repository..."
curl -o docker-compose.yml ${DOWNLOAD_URL_DOCKER_COMPOSE}
curl -o gosharexserver-config.toml ${DOWNLOAD_URL_DOCKER_COMPOSE_CONFIG}
echo "Download complete. Creating MongoDB data directory..."
mkdir ./data
echo "Done with installation. Use \"docker-compose up\" to start."
