APP_NAME := contract-cli
OUTPUT := build
VERSION ?= 0.0.1

default: test build

test: 
	go test ./... -v

build:
	go build -o ${OUTPUT}/${APP_NAME}

release:
	GOOS=linux   GOARCH=amd64   go build -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_amd64
	GOOS=linux   GOARCH=arm64   go build -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_arm64
	GOOS=linux   GOARCH=s390x   go build -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_s390x
	GOOS=linux   GOARCH=ppc64le go build -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_ppc64le
	GOOS=darwin  GOARCH=amd64   go build -o ${OUTPUT}/${APP_NAME}_${VERSION}_darwin_amd64
	GOOS=darwin  GOARCH=arm64   go build -o ${OUTPUT}/${APP_NAME}_${VERSION}_darwin_arm64
	GOOS=windows GOARCH=amd64   go build -o ${OUTPUT}/${APP_NAME}_${VERSION}_windows_amd64
	GOOS=windows GOARCH=arm64   go build -o ${OUTPUT}/${APP_NAME}_${VERSION}_windows_arm64

test-cover:
	go test -coverprofile build/cover.out ./...
	go tool cover -html=build/cover.out

update-packages:
	go get -u all

tidy:
	go mod tidy

clean:
	find ./build ! -name '.gitkeep' -type f -delete
