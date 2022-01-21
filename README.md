# Multimedia encodings referencing endpoints

This is a re-implementation of the encodings endpoints, based on Go and API Gateway instead of PHP.

## What's in the box?

Main pieces:
- **referenceapi/** - the `reference` endpoint. This looks up the content and gives the result as a single line of text
- **migration/** - a commandline tool (NOT a lambda function!) to migrate data from MySQL into DynamoDB

Other bits:
- common/ - functionality that is shared between all the endpoints. This is the code that does the actual database scanning
- infra/  - Cloudformation deployments to make it all work

## Development Prerequisites

You will need:
- A recent version of Go, at least 1.11 but preferably 1.16 or 1.17
- GNU make - this is normally available as standard on Mac or Linux and can be installed (usually via Cygwin) on Windows

To make sure everything is working, go to the root of this repository in your terminal and run:
```bash
make test
```

## First-time setup

The process of performing a first-time deployment must be done in the right order, otherwise steps will fail because
dependencies from previous steps are not there

0. Make sure you have "Development Prerequisites", above
1. Make sure you have a writable bucket to put the compiled lambda functions into
2. Decide on the `App`, `Stack` and `Stage` identifiers you will use (`Stage` must be `PROD` or `CODE`)
3. In a terminal, set all of these as environment variables:
```bash
declare -x DEPLOYMENTBUCKET=your-deployment-bucket
declare -x APP=encodings-endpoint
declare -x STACK=multimedia
declare -x UPLOADVERSION=main
```
5. Run `make upload` from the root of this repository. This will compile and upload the lambda function code
6. In the AWS Web console, go to Cloudformation and Create Stack
7. Use the file `infra/apigateway_base.yaml`.  When deploying, make a note of the name you choose; you'll need it later.
Make sure you use the same `App` and `Stack` tags that you used earlier.  This will set up the basic, shared API Gateway setup.
8. With `apigateway_base.yaml` set up, then deploy `infra/endpoints.yaml`.  Use the stack name that you chose for
`apigateway_base` as the `APIGatewayStack` parameter and make sure that the deployment bucket, app, stack and stage
parameters are EXACTLY the ones you used for upload (these values are used to compute the path to the code bundle)

With this in place, you can go to the API Gateway app in the AWS Console and test the API that way. You can also
retrieve a "direct access" url.

## Development updates

Once that has been done successfully, you can deploy updates of the lambda functions directly with the `make` utility,
provided you have the environment variables set:

NOTE - about versioning and deployment - TODO
```bash
declare -x DEPLOYMENTBUCKET=your-deployment-bucket
declare -x APP=encodings-endpoint
declare -x STACK=multimedia
declare -x STAGE=CODE
make clean && make deploy
```

## Proper DNS names

TODO
