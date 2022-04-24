# Copyright (C) 2022 xaxys. All rights reserved.
PACKAGE_NAME          := maintainman

ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
	GO		   ?= go.exe
    PWD 	   := ${CURDIR}
	TARGET	   := $(PACKAGE_NAME).exe
	BUILD_TAGS := $(shell git describe --tags --always --dirty="-dev")
	BUILD_TIME := $(shell echo %date% %time%)
	GIT_COMMIT := $(shell git rev-parse --short HEAD)
	GO_VERSION := $(shell go version)
	GOPATH     := $(subst ;,,$(shell go env GOPATH))
	RM_CMD_1   := del /s /q
	RM_CMD_2   := 
	EXPORT     := set
else
	GO		   ?= go
    PWD 	   := ${CURDIR}
	TARGET	   := $(PACKAGE_NAME)
	BUILD_TAGS := $(shell git describe --tags --always --dirty="-dev")
	BUILD_TIME := $(shell date --utc)
	GIT_COMMIT := $(shell git rev-parse --short HEAD)
	GO_VERSION := $(shell go version)
	GOPATH     := $(shell go env GOPATH)
	RM_CMD_1   := find . -type f -name
	RM_CMD_2   := -delete
	EXPORT     := export
endif

all: build

build:
	@echo "Building MaintainMan ..."
	@$(GO) env -w CGO_ENABLED="1"
	@$(GO) build \
		-ldflags="-X 'main.BuildTags=$(BUILD_TAGS)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GoVersion=$(GO_VERSION)'" \
		-o $(TARGET) $(PWD)/main.go

test: test-short

test-full: clean
	@echo "Testing MaintainMan ..."
	@$(GO) env -w CGO_ENABLED="1"
	@$(GO) test \
		-ldflags="-X 'main.BuildTags=$(BUILD_TAGS)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GoVersion=$(GO_VERSION)'" \
		-timeout=30m -coverprofile=coverage.out ./...

test-short: clean
	@echo "Testing MaintainMan ..."
	@$(GO) env -w CGO_ENABLED="1"
	@$(GO) test \
		-ldflags="-X 'main.BuildTags=$(BUILD_TAGS)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GoVersion=$(GO_VERSION)'" \
		-timeout=30m -short -coverprofile=coverage.out ./...

clean:
	@echo "Cleaning MaintainMan ..."
	@$(RM_CMD_1) $(TARGET)    $(RM_CMD_2)
	@$(RM_CMD_1) coverage.out $(RM_CMD_2)
	@$(RM_CMD_1) "*.db"       $(RM_CMD_2)
	@$(RM_CMD_1) "*.exe"      $(RM_CMD_2)
	@$(RM_CMD_1) "*.out"      $(RM_CMD_2)
	@$(RM_CMD_1) "*.yaml"     $(RM_CMD_2)

.PHONY: all test test-full test-short clean