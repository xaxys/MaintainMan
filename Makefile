# Copyright (C) 2019-2020 xaxys. All rights reserved.

ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
	GO		?= go.exe
    PWD 	:= ${CURDIR}
	TARGET	:= maintainman.exe
	BUILD_TAGS := $(shell git describe --tags --always --dirty="-dev")
	BUILD_TIME := $(shell echo %date% %time%)
	GIT_COMMIT := $(shell git rev-parse --short HEAD)
	GO_VERSION := $(shell go version)
else
	GO		?= go
    PWD 	:= ${CURDIR}
	TARGET	:= maintainman
	BUILD_TAGS := $(shell git describe --tags --always --dirty="-dev")
	BUILD_TIME := $(shell date --utc)
	GIT_COMMIT := $(shell git rev-parse --short HEAD)
	GO_VERSION := $(shell go version)
endif

all: build

build:
	@echo "Building MaintainMan ..."
	@$(GO) env -w CGO_ENABLED="1"
	@$(GO) build \
		-ldflags="-X 'main.BuildTags=$(BUILD_TAGS)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GoVersion=$(GO_VERSION)'" \
		-o $(TARGET) $(PWD)/main.go

test:
	@echo "Testing MaintainMan ..."
	@$(GO) env -w CGO_ENABLED="1"
	@$(GO) test \
		-ldflags="-X 'main.BuildTags=$(BUILD_TAGS)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GoVersion=$(GO_VERSION)'" \
		-coverprofile=coverage.out ./...
