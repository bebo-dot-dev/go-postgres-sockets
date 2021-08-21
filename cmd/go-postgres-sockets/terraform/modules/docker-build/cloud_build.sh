#!/usr/bin/env bash

#
# A bash shell script that performs a GCP gcloud call to push a container into the terraform created GCP artefact repository
#
# This shell script is called by local-exec in ./modules/docker-build terraform module main.tf config and
# required variables are passed to this script as positional arguments
#

REGION=$1;
PROJECT_ID=$2;
REPOSITORY_ID=$3;
CONTAINER_ID=$4;

# get the working directory of this script to correctly workout the location of the application Docker file
THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd $THIS_SCRIPT_DIR/../../../

gcloud auth configure-docker $REGION-docker.pkg.dev
gcloud builds submit --tag $REGION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY_ID/$CONTAINER_ID .