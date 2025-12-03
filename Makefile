APP_NAME := contract-cli
OUTPUT := build
VERSION ?= 0.0.1
BUILD_DATE ?= $(shell date -u)

LDFLAGS_COMMON := -s -w -X 'main.version=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'
LDFLAGS_DEV := -X 'main.version=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'

default: test

help:
	@echo "Available targets:"
	@echo "  make build           - Build binary for local development (Linux amd64)"
	@echo "  make release         - Build binaries for all supported platforms"
	@echo "  make test            - Run all tests"
	@echo "  make test-cover      - Run tests with coverage report"
	@echo "  make fmt             - Format Go code"
	@echo "  make tidy            - Tidy Go modules"
	@echo "  make clean           - Remove build artifacts"
	@echo "  make install-deps    - Download Go module dependencies"
	@echo "  make update-packages - Update all Go packages"

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
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${OUTPUT}/${APP_NAME} -ldflags "${LDFLAGS_DEV} -X 'main.osName=Linux' -X 'main.osArch=x86_64'"

release:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64   go build -ldflags "${LDFLAGS_COMMON} -X 'main.osName=Linux'   -X 'main.osArch=x86_64'"  -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_amd64
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64   go build -ldflags "${LDFLAGS_COMMON} -X 'main.osName=Linux'   -X 'main.osArch=ARM64'"   -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_arm64
	CGO_ENABLED=0 GOOS=linux   GOARCH=s390x   go build -ldflags "${LDFLAGS_COMMON} -X 'main.osName=Linux'   -X 'main.osArch=S390x'"   -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_s390x
	CGO_ENABLED=0 GOOS=linux   GOARCH=ppc64le go build -ldflags "${LDFLAGS_COMMON} -X 'main.osName=Linux'   -X 'main.osArch=Ppc64le'" -o ${OUTPUT}/${APP_NAME}_${VERSION}_linux_ppc64le
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64   go build -ldflags "${LDFLAGS_COMMON} -X 'main.osName=Darwin'  -X 'main.osArch=x86_64'"  -o ${OUTPUT}/${APP_NAME}_${VERSION}_darwin_amd64
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64   go build -ldflags "${LDFLAGS_COMMON} -X 'main.osName=Darwin'  -X 'main.osArch=ARM64'"   -o ${OUTPUT}/${APP_NAME}_${VERSION}_darwin_arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64   go build -ldflags "${LDFLAGS_COMMON} -X 'main.osName=Windows' -X 'main.osArch=x86_64'"  -o ${OUTPUT}/${APP_NAME}_${VERSION}_windows_amd64.exe
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64   go build -ldflags "${LDFLAGS_COMMON} -X 'main.osName=Windows' -X 'main.osArch=ARM64'"   -o ${OUTPUT}/${APP_NAME}_${VERSION}_windows_arm64.exe

.PHONY: default help install-deps test build release test-cover update-packages tidy clean fmt
