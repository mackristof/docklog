#!/bin/bash
go build -o tick tick.go
docker build -t tick .
# docker kill toto1 toto2
docker run -d --rm -l speaker=toto --name toto1 tick 
docker run -d --rm -l speaker=tata --name toto2 tick 
