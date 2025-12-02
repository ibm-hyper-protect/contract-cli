APP_NAME := contract-cli
OUTPUT := build
VERSION ?= 0.0.1
BUILD_DATE ?= $(shell date -u)

default: test build

install-deps:
	go mod download

test:
	go test ./... -v

test-cover:
	go test -coverprofile build/cover.out ./...
	go tool cover -html=build/cover.out

update-packages:
	go get -u all

tidy:
	go mod tidy

clean:
	find ./build ! -name '.gitkeep' -type f -delete

fmt:
	go fmt ./...

build:
	GOOS=linux GOARCH=amd64 go build -o ${OUTPUT}/${APP_NAME} -ldflags "-X 'main.version=${VERSION}' -X 'main.osName=Linux' -X 'main.osArch=x86_64' -X 'main.buildDate=${BUILD_DATE}'"

release:
	GOOS=linux   GOARCH=amd64   go build -ldflags "-s -w -X 'main.version=${VERSION}' -X 'main.osName=Linux'   -X 'main.osArch=x86_64'  -X 'main.buildDate=${BUILD_DATE}'" -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_amd64
	GOOS=linux   GOARCH=arm64   go build -ldflags "-s -w -X 'main.version=${VERSION}' -X 'main.osName=Linux'   -X 'main.osArch=ARM64'   -X 'main.buildDate=${BUILD_DATE}'" -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_arm64
	GOOS=linux   GOARCH=s390x   go build -ldflags "-s -w -X 'main.version=${VERSION}' -X 'main.osName=Linux'   -X 'main.osArch=S390x'   -X 'main.buildDate=${BUILD_DATE}'" -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_s390x
	GOOS=linux   GOARCH=ppc64le go build -ldflags "-s -w -X 'main.version=${VERSION}' -X 'main.osName=Linux'   -X 'main.osArch=Ppc64le' -X 'main.buildDate=${BUILD_DATE}'" -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_ppc64le
	GOOS=darwin  GOARCH=amd64   go build -ldflags "-s -w -X 'main.version=${VERSION}' -X 'main.osName=Darwin'  -X 'main.osArch=x86_64'  -X 'main.buildDate=${BUILD_DATE}'" -o ${OUTPUT}/${APP_NAME}_${VERSION}_darwin_amd64
	GOOS=darwin  GOARCH=arm64   go build -ldflags "-s -w -X 'main.version=${VERSION}' -X 'main.osName=Darwin'  -X 'main.osArch=ARM64'   -X 'main.buildDate=${BUILD_DATE}'" -o ${OUTPUT}/${APP_NAME}_${VERSION}_darwin_arm64
	GOOS=windows GOARCH=amd64   go build -ldflags "-s -w -X 'main.version=${VERSION}' -X 'main.osName=Windows' -X 'main.osArch=x86_64'  -X 'main.buildDate=${BUILD_DATE}'" -o ${OUTPUT}/${APP_NAME}_${VERSION}_windows_amd64.exe
	GOOS=windows GOARCH=arm64   go build -ldflags "-s -w -X 'main.version=${VERSION}' -X 'main.osName=Windows' -X 'main.osArch=ARM64'   -X 'main.buildDate=${BUILD_DATE}'" -o ${OUTPUT}/${APP_NAME}_${VERSION}_windows_arm64.exe

.PHONY: default install-deps test build release test-cover update-packages tidy clean fmt
