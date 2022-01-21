#!/bin/bash -e

##Expected arguments:
# $1 - lambda function name
# $2 - alias to deploy it with

PUBLISHED=0
while [[ "$PUBLISHED" == "0" ]]; do
  sleep 2s
  aws lambda publish-version --function-name "$1" > published-version.json
  if [[ "$?" == "0" ]]; then PUBLISHED=1; fi
done

VERS=$(jq .Version < published-version.json | sed s/\"//g) #extract the version number with JQ. It comes as a string so we must strip out the quotes.
echo "We have just (re-)deployed version $VERS"
aws lambda update-alias --function-name "$1" --name "$2" --function-version ${VERS}