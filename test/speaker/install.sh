#!/bin/bash
go build -o tick main.go
docker build -t tick .