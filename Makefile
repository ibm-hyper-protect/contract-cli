TOOL_NAME=contract-cli

default: test build

build:
	go build -o ${TOOL_NAME} .

test: 
	go test ./...

test-cover:
	go test -coverprofile build/cover.out ./...
	go tool cover -html=build/cover.out

update-packages:
	go get -u all

tidy:
	go mod tidy
