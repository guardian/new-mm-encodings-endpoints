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

clean:
	make -C migration/ clean
	make -C referenceapi/ clean

deploy:
	make -C referenceapi/ deploy
