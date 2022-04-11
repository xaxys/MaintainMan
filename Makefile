# Copyright (C) 2022 xaxys. All rights reserved.
PACKAGE_NAME          := maintainman

ifeq ($(OS),Windows_NT)  # is Windows_NT on XP, 2000, 7, Vista, 10...
	GO		?= go.exe
    PWD 	:= ${CURDIR}
	TARGET	:= $(PACKAGE_NAME).exe
	BUILD_TAGS := $(shell git describe --tags --always --dirty="-dev")
	BUILD_TIME := $(shell echo %date% %time%)
	GIT_COMMIT := $(shell git rev-parse --short HEAD)
	GO_VERSION := $(shell go version)
	RM := del /s /q
else
	GO		?= go
    PWD 	:= ${CURDIR}
	TARGET	:= $(PACKAGE_NAME)
	BUILD_TAGS := $(shell git describe --tags --always --dirty="-dev")
	BUILD_TIME := $(shell date --utc)
	GIT_COMMIT := $(shell git rev-parse --short HEAD)
	GO_VERSION := $(shell go version)
	RM := rm -rf
endif

all: build

build: bindata
	@echo "Building MaintainMan ..."
	@$(GO) env -w CGO_ENABLED="1"
	@$(GO) build \
		-ldflags="-X 'main.BuildTags=$(BUILD_TAGS)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GoVersion=$(GO_VERSION)'" \
		-o $(TARGET) $(PWD)/main.go

test: clean bindata
	@echo "Testing MaintainMan ..."
	@$(GO) env -w CGO_ENABLED="1"
	@$(GO) test \
		-ldflags="-X 'main.BuildTags=$(BUILD_TAGS)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GoVersion=$(GO_VERSION)'" \
		-coverprofile=coverage.out

bindata:
	@echo "Run go-bindata ..."
	go-bindata -nomemcopy --pkg bindata -o ./bindata/bindata.go fonts/...

clean:
	@echo "Cleaning MaintainMan ..."
	@$(RM) $(TARGET)
	@$(RM) coverage.out
	@$(RM) *.db
	@$(RM) *.exe
	@$(RM) *.out
	@$(RM) *.yaml

.PHONY: all test bindata clean