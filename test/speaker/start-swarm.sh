#!/bin/bash
# set -o errexit
set -o pipefail

if hash docker-machine 2>/dev/null; then
        echo "docker-machine found no need to install"
    else

        # install virtualbox
        sudo apt-get install virtualbox

        # install docker-machine
        curl -L https://github.com/docker/machine/releases/download/v0.10.0/docker-machine-`uname -s`-`uname -m` >/tmp/docker-machine && chmod +x /tmp/docker-machine && sudo cp /tmp/docker-machine /usr/local/bin/docker-machine
fi
docker-machine rm -y -f swarm-manager
docker-machine rm -y -f worker-swarm
docker-machine create --engine-env 'DOCKER_OPTS="-H unix:///var/run/docker.sock"' --driver virtualbox swarm-manager
docker-machine create --engine-env 'DOCKER_OPTS="-H unix:///var/run/docker.sock"' --driver virtualbox worker-swarm

ip_swarm=$(docker-machine ip swarm-manager)
export ip_swarm
eval "$(docker-machine env swarm-manager)"
docker swarm init --listen-addr $ip_swarm --advertise-addr $ip_swarm


# getting token from swarm leader
token=$(docker swarm join-token worker -q)
#init swarm slave worker-swarm
eval "$(docker-machine env worker-swarm)"
docker swarm join --token $token $ip_swarm:2377

eval "$(docker-machine env swarm-manager)"
go build -o tick tick.go
docker build -t tick .
docker service rm toto1 toto2
docker service create --detach=false -l speaker=toto --name toto1 tick 
docker service create --detach=false -l speaker=tata --name toto2 tick 
