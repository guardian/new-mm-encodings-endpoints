all: test-against-captureddata.linuxx64 test-against-captureddata.linuxarm test-against-captureddata.macx64 test-against-captureddata.macarm

clean:
	rm -f test-against-captureddata.linux* test-against-captureddata.mac* test-against-captureddata

test-against-captureddata.macx64: main.go async_reader.go async_tester.go EndpointEvent.go
	GOOS=darwin GOARCH=amd64 go build -o test-against-captureddata.macx64

test-against-captureddata.macarm: main.go async_reader.go async_tester.go EndpointEvent.go
	GOOS=darwin GOARCH=arm64 go build -o test-against-captureddata.macarm

test-against-captureddata.linuxx64: main.go async_reader.go async_tester.go EndpointEvent.go
	GOOS=linux GOARCH=amd64 go build -o test-against-captureddata.linuxx64

test-against-captureddata.linuxarm: main.go async_reader.go async_tester.go EndpointEvent.go
	GOOS=linux GOARCH=arm64 go build -o test-against-captureddata.linuxarm

