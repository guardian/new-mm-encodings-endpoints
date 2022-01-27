.PHONY: referenceapi upload clean deploy migration

all: referenceapi migration

referenceapi:
	make -C referenceapi/

upload:
	make -C referenceapi/ upload

migration:
	make -C migration/

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

deploy:
	make -C referenceapi/ deploy
