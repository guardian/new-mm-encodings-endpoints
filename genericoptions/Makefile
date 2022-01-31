.PHONY: all

all: genericoptions.zip

genericoptions: genericoptions.go ../common/responses.go
	GOOS=linux GOARCH=amd64 go build -o genericoptions

genericoptions.zip: genericoptions
	zip genericoptions.zip genericoptions

upload: genericoptions.zip
	../ci-scripts/upload-and-deploy.sh "genericoptions.zip"

deploy: genericoptions.zip
	../ci-scripts/upload-and-deploy.sh "genericoptions.zip" "${APP}-GenericOptions"

clean:
	rm -f genericoptions genericoptions.zip published-version.json