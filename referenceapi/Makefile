.PHONY: all

all: referenceapi.zip

referenceapi: referenceapi.go ../common/config.go ../common/find_content.go ../common/idmapping.go ../common/responses.go
	GOOS=linux GOARCH=amd64 go build -o referenceapi

referenceapi.zip: referenceapi
	zip referenceapi.zip referenceapi

upload: referenceapi.zip
	../ci-scripts/upload-and-deploy.sh "referenceapi.zip"

deploy: referenceapi.zip
	../ci-scripts/upload-and-deploy.sh "referenceapi.zip" "${APP}-References"

clean:
	rm -f referenceapi referenceapi.zip