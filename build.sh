#!/bin/bash
if [ -z $1 ];
then
    echo "Usage: build.sh <docker user> <version>"
    echo "This builds and pushes an updated image to a docker repo"
    exit 1
fi
docker build . -t $1/go-scrape:$2
docker push $1/go-scrape:$2
