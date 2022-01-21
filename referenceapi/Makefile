.PHONY: all

all: referenceapi.zip

referenceapi: referenceapi.go ../common/config.go ../common/find_content.go ../common/idmapping.go ../common/responses.go
	GOOS=linux GOARCH=amd64 go build -o referenceapi

referenceapi.zip: referenceapi
	zip referenceapi.zip referenceapi

upload: referenceapi.zip
	aws s3 cp referenceapi.zip s3://${DEPLOYMENTBUCKET}/${APP}/${STACK}/${STAGE}/referenceapi.zip

deploy: upload

clean:
	rm -f referenceapi referenceapi.zip