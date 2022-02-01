.PHONY: referenceapi genericoptions upload clean deploy migration test-against-captureddata video mediatag

all: referenceapi genericoptions migration test-against-captureddata video mediatag

referenceapi:
	make -C referenceapi/

upload:
	make -C referenceapi/ upload
	make -C genericoptions/ upload
	make -C video/ upload
	make -C mediatag/ upload

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
	make -C mediatag/ clean

deploy:
	make -C referenceapi/ deploy
	make -C genericoptions/ deploy
	make -C video/ deploy
	make -C mediatag/ deploy

video:
	make -C video/

mediatag:
	make -C mediatag/
