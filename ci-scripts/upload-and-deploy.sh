#!/bin/bash

#This script expects the following arguments:
# $1 - the zipfile to upload
# $2 - the function name to associate it with

if [ ! -x "`which jq`" ]; then
  echo This script requires the \"jq\" utility. Please install it and re-run.
  exit 1
fi

if [ ! -x "`which aws`" ]; then
  echo This script requires the \"aws\" commandline utility. Please install it and re-run.
  exit 1
fi

if [ "${DEPLOYMENTBUCKET}" == "" ]; then
  echo You must set DEPLOYMENTBUCKET before uploading
  exit 1
fi

if [ "${APP}" == "" ]; then
  echo You must  set APP before uploading
  exit 1
fi

if [ "${STACK}" == "" ]; then
  echo You must set STACK before uploading
  exit 1
fi

if [ "${UPLOADVERSION}" == "" ]; then
  UPLOADVERSION=$(uuidgen)
fi

echo Uploading to "s3://${DEPLOYMENTBUCKET}/${APP}/${STACK}/${UPLOADVERSION}/$1"

aws s3 cp "$1" "s3://${DEPLOYMENTBUCKET}/${APP}/${STACK}/${UPLOADVERSION}/$1"

if [ "$?" != "0" ]; then
  echo Upload failed!
  exit 1
fi

#if there is no function name specified then exit now
if [ "$2" == "" ]; then
  exit 0
fi

aws lambda update-function-code --function-name "$2" --s3-bucket "${DEPLOYMENTBUCKET}" --s3-key "${APP}/${STACK}/${UPLOADVERSION}/$1" > /dev/null
if [ "$?" != "0" ]; then
  echo Could not post update to lambda!
  exit 1
fi

PUBLISHED=0
CTR=0
while [[ "$PUBLISHED" == "0" ]]; do
  sleep 2s
  aws lambda publish-version --function-name "$2" > published-version.json
  if [[ "$?" == "0" ]]; then PUBLISHED=1; fi
  CTR=$((CTR+1))
  if [[ "$CTR" == "10" ]]; then
    echo Could not publish after 10 attempts, giving up
    exit 1
  fi
done

VERS=$(jq .Version < published-version.json | sed s/\"//g) #extract the version number with JQ. It comes as a string so we must strip out the quotes.
echo "We have just (re-)deployed version $VERS"

if [ "${GITHUB_REF_NAME}" != "" ]; then
  echo Running on Github from branch ${GITHUB_REF_NAME}
  aws lambda create-alias --function-name "$2" --name "${GITHUB_REF_NAME}" --function-version "${VERS}"
  if [ "$?" != "0" ]; then
    echo Could not create function alias for version ${VERS} to ${GITHUB_REF_NAME}
    exit 1
  fi
fi

if [ "${STAGE}" == "" ]; then
  echo Not linking to any deployment as STAGE is not set
else
  aws lambda update-alias --function-name "$2" --name "${STAGE}" --function-version "${VERS}"
fi