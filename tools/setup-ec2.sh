#!/usr/bin/env bash

set -eux
set -o pipefail

sudo apt get-update
sudo apt-get install docker.io -y
sudo systemctl start-docker
docker ps
sudo chmod 666 /var/run/docker-sock
sudo systemctl enable docker
docker --version