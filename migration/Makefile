all: migration.linuxx64 migration.linuxarm migration.macx64 migration.macarm

clean:
	rm -f migration.linux* migration.mac* migration

migration.macx64: main.go dynamo_writer.go async_db_reader.go
	GOOS=darwin GOARCH=amd64 go build -o migration.macx64

migration.macarm: main.go dynamo_writer.go async_db_reader.go
	GOOS=darwin GOARCH=arm64 go build -o migration.macarm

migration.linuxx64: main.go dynamo_writer.go async_db_reader.go
	GOOS=linux GOARCH=amd64 go build -o migration.linuxx64

migration.linuxarm: main.go dynamo_writer.go async_db_reader.go
	GOOS=linux GOARCH=arm64 go build -o migration.linuxarm

