#!/bin/sh

cd angular

docker login -u $CI_REGISTRY_USER -p $CI_JOB_TOKEN $CI_REGISTRY
VERSION=$CI_COMMIT_SHORT_SHA
docker build -t $REGISTRY_URL/$PROJECT_GROUP/frontend:$VERSION -t $REGISTRY_URL/$PROJECT_GROUP/frontend:latest -f .docker/Dockerfile .
docker push $REGISTRY_URL/$PROJECT_GROUP/frontend:$VERSION
docker push $REGISTRY_URL/$PROJECT_GROUP/frontend:latest
