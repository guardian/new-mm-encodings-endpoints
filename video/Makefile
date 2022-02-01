.PHONY: all

all: video.zip

video: video.go ../common/config.go ../common/find_content.go ../common/idmapping.go ../common/responses.go
	GOOS=linux GOARCH=amd64 go build -o video

video.zip: video
	zip video.zip video

upload: video.zip
	../ci-scripts/upload-and-deploy.sh "video.zip"

deploy: video.zip
	../ci-scripts/upload-and-deploy.sh "video.zip" "${APP}-Video"

clean:
	rm -f video video.zip published-version.json