# Copyright (c) 2025 IBM Corp.
# All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

APP_NAME := contract-cli
OUTPUT := build
VERSION ?= 0.0.0
BUILD_DATE ?= $(shell date -u)

# Auto-detect host platform
HOST_OS := $(shell go env GOOS)
HOST_ARCH := $(shell go env GOARCH)

# Map GOOS to display name for ldflags
OS_NAME := $(if $(filter linux,$(HOST_OS)),Linux,$(if $(filter darwin,$(HOST_OS)),Darwin,Windows))

# Map GOARCH to display name for ldflags
OS_ARCH := $(if $(filter amd64,$(HOST_ARCH)),x86_64,$(if $(filter arm64,$(HOST_ARCH)),ARM64,$(if $(filter s390x,$(HOST_ARCH)),S390x,Ppc64le)))

LDFLAGS_DEV := -X 'main.version=${VERSION}' -X 'main.buildDate=${BUILD_DATE}' -X 'main.osName=${OS_NAME}' -X 'main.osArch=${OS_ARCH}'

default: test

help:
	@echo "Available targets:"
	@echo "  make build           - Build binary for local development (auto-detects host OS/arch)"
	@echo "  make snapshot        - Build all artifacts locally via GoReleaser (no publish)"
	@echo "  make release-local   - Full release dry-run via GoReleaser (no publish)"
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
	CGO_ENABLED=0 go build -o ${OUTPUT}/${APP_NAME} -ldflags "${LDFLAGS_DEV}"

snapshot:
	CHANGELOG_DISABLE=true goreleaser build --snapshot --clean

release-local:
	CHANGELOG_DISABLE=true goreleaser release --snapshot --clean

.PHONY: default help install-deps test build snapshot release-local test-cover update-packages tidy clean fmt
