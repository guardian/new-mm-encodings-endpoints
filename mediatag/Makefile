.PHONY: all

all: mediatag.zip

mediatag: mediatag.go ../common/config.go ../common/find_content.go ../common/idmapping.go ../common/responses.go
	GOOS=linux GOARCH=amd64 go build -o mediatag

mediatag.zip: mediatag
	zip mediatag.zip mediatag

upload: mediatag.zip
	../ci-scripts/upload-and-deploy.sh "mediatag.zip"

deploy: mediatag.zip
	../ci-scripts/upload-and-deploy.sh "mediatag.zip" "${APP}-MediaTag"

clean:
	rm -f mediatag mediatag.zip published-version.json