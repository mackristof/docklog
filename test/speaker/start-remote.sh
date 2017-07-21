#!/bin/bash
set -o errexit
set -o pipefail

if hash docker-machine 2>/dev/null; then
        echo "docker-machine found no need to install"
    else

        # install virtualbox
        sudo apt-get install virtualbox

        # install docker-machine
        curl -L https://github.com/docker/machine/releases/download/v0.10.0/docker-machine-`uname -s`-`uname -m` >/tmp/docker-machine && chmod +x /tmp/docker-machine && sudo cp /tmp/docker-machine /usr/local/bin/docker-machine
fi

# docker-machine create --engine-env 'DOCKER_OPTS="-H unix:///var/run/docker.sock"' --driver virtualbox --virtualbox-memory "1024" leader1
ip_leader1=$(docker-machine ip leader1)
export ip_leader1

eval "$(docker-machine env leader1)"
go build -o tick tick.go
docker build -t tick .
docker kill toto1 toto2
docker run -d --rm -l speaker=toto --name toto1 tick 
docker run -d --rm -l speaker=tata --name toto2 tick 
