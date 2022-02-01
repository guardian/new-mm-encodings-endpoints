.PHONY: referenceapi genericoptions upload clean deploy migration test-against-captureddata video

all: referenceapi genericoptions migration test-against-captureddata video

referenceapi:
	make -C referenceapi/

upload:
	make -C referenceapi/ upload
	make -C genericoptions/ upload
	make -C video/ upload

migration:
	make -C migration/


genericoptions:
	make -C genericoptions/

test-against-captureddata:
	make -C test-against-captureddata/

test:
	go test ./...

coverage:
	rm -f cover.out
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out
	rm -f cover.out

clean:
	rm -f cover.out
	make -C migration/ clean
	make -C referenceapi/ clean
	make -C genericoptions/ clean
	make -C test-against-captureddata/ clean
	make -C video/ clean

deploy:
	make -C referenceapi/ deploy
	make -C genericoptions/ deploy
	make -C video/ deploy

video:
	make -C video/
