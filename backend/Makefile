# indece Monitor
# Copyright (C) 2023 indece UG (haftungsbeschränkt)
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License or any
# later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <https://www.gnu.org/licenses/>.

PROJECT_NAME ?= $(shell basename $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
GOCMD=go
GOPATH=$(shell $(GOCMD) env GOPATH))
GOBUILD=$(GOCMD) build
GOGENERATE=$(GOCMD) generate
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
DIR_SOURCE=./src
DIR_DIST=./dist
DIR_GENERATED_MODEL=$(DIR_SOURCE)/generated/model
BINARY_NAME=$(DIR_DIST)/bin/$(PROJECT_NAME)
BUILD_DATE=$(shell date +%Y%m%d.%H%M%S)
BUILD_VERSION ?= $(shell git rev-parse --short HEAD)
LDFLAGS := 
LDFLAGS := $(LDFLAGS) -X github.com/indece-official/monitor/backend/src/buildvars.ProjectName=$(PROJECT_NAME)
LDFLAGS := $(LDFLAGS) -X github.com/indece-official/monitor/backend/src/buildvars.BuildDate=$(BUILD_DATE)
LDFLAGS := $(LDFLAGS) -X github.com/indece-official/monitor/backend/src/buildvars.BuildVersion=$(BUILD_VERSION)

all: generate test sec build

generate:
	mkdir -p $(DIR_GENERATED_MODEL)
	rm -rf $(DIR_GENERATED_MODEL)/*
	$(GOGENERATE) -tags=bindata ./...

sec:
	gosec -conf ./gosec.json ./...

build:
	mkdir -p $(DIR_DIST)/bin
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) -tags=prod -v $(DIR_SOURCE)/main.go

test:
	mkdir -p $(DIR_DIST)
ifeq ($(OUTPUT),json)
	$(GOTEST) -v ./...  -cover -coverprofile $(DIR_DIST)/cover.out -json > $(DIR_DIST)/test.json
else
	$(GOTEST) -v ./...  -cover
endif

clean:
	#$(GOCLEAN)
	rm -rf $(DIR_DIST)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	echo test
	#$(GOGET) -d -v ./...
