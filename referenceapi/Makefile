.PHONY: all

all: referenceapi.zip

referenceapi: referenceapi.go ../common/config.go ../common/find_content.go ../common/idmapping.go ../common/responses.go
	GOOS=linux GOARCH=amd64 go build -o referenceapi

referenceapi.zip: referenceapi
	zip referenceapi.zip referenceapi

upload: referenceapi.zip
	. ../ci-scripts/set_upload_version.sh
	aws s3 cp referenceapi.zip s3://${DEPLOYMENTBUCKET}/${APP}/${STACK}/${UPLOADVERSION}/referenceapi.zip

deploy: upload
	. ../ci-scripts/set_upload_version.sh
	aws lambda update-function-code --function-name ${APP}-References --s3-bucket ${DEPLOYMENTBUCKET} --s3-key ${APP}/${STACK}/${UPLOADVERSION}/referenceapi.zip
	../ci-scripts/lambda-version.sh "${APP}-References" "${STAGE}"

clean:
	rm -f referenceapi referenceapi.zip