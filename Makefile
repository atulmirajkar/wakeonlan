BINARY_NAME=wakeonlan

all: build test

build:
	go build -o ${BINARY_NAME} wol.go

test:
	go test .

clean:
	go clean
	rm ${BINARY_NAME}