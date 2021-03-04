#!/bin/bash

set -e

if [ -z "$ORGANIZATION" ]; then
  echo "ORGANIZATION is not set. Quitting."
  exit 1
fi

if [ -z "$PROJECT" ]; then
  echo "PROJECT is not set. Quitting."
  exit 1
fi

if [ -z "$SERVICE" ]; then
  echo "SERVICE is not set. Quitting."
  exit 1
fi

if [ -z "$DOCKER_USER" ]; then
  echo "DOCKER_USER is not set. Quitting."
  exit 1
fi

if [ -z "$DOCKER_TOKEN" ]; then
  echo "DOCKER_TOKEN is not set. Quitting."
  exit 1
fi

if [ -z "$DOCKER_REPOSITORY" ]; then
  echo "DOCKER_REPOSITORY is not set. Quitting."
  exit 1
fi

if [ -z "$DOCKER_REPOSITORY_TAG" ]; then
  echo "DOCKER_REPOSITORY_TAG is not set. Quitting."
  exit 1
fi

if [ -z "$IAAS_CRED" ]; then
  echo "IAAS_CRED is not set. Quitting."
  exit 1
fi

if [ -z "$IAAS_PROVIDER" ]; then
  echo "IAAS_PROVIDER is not set. Quitting."
  exit 1
fi

if [ -z "$IAAS_REGION" ]; then
  echo "IAAS_REGION is not set. Quitting."
  exit 1
fi

if [ -z "$NODE_TYPE" ]; then
  echo "NODE_TYPE is not set. Quitting."
  exit 1
fi

echo "Create project if not exists"
if ! /sqsc project get -project-name $ORGANIZATION}/$PROJECT; then
  /sqsc project create \
    -credential $IAAS_CRED \
    -monitoring netdata \
    -name $PROJECT \
    -node-size $NODE_TYPE \
    -infra-type high-availability \
    -organization $ORGANIZATION \
    -provider $IAAS_PROVIDER \
    -region $IAAS_REGION \
    -yes
fi

echo "Create service if not exists"
if ! /sqsc container list --project-name $ORGANIZATION}/$PROJECT | grep $SERVICE; then
  /sqsc container add \
    -project-name $ORGANIZATION/$PROJECT \
    -servicename $SERVICE \
    -name $DOCKER_REPOSITORY:$DOCKER_REPOSITORY_TAG \
    -username $DOCKER_USER \
    -password $DOCKER_TOKEN
fi

echo "Open HTTP port"
NETWORK_RULES_NAME=http
if ! /sqsc network-rule list -project-name $ORGANIZATION/$PROJECT -service-name $SERVICE | grep $NETWORK_RULES_NAME; then
  /sqsc \
    network-rule create \
    -project-name $ORGANIZATION/$PROJECT \
    -external-protocol http \
    -internal-port 80 \
    -internal-protocol http \
    -name $NETWORK_RULES_NAME \
    -service-name $SERVICE
fi

echo "Schedule web service"
/sqsc service schedule --project-name $ORGANIZATION/$PROJECT $SERVICE