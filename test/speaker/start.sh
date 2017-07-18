#!/bin/bash
docker run -d -l speaker=toto --name toto1 tick 
docker run -d -l speaker=tata --name toto2 tick 
