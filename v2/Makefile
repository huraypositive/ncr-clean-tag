COMMAND=nct
GIT_VER=$(shell git describe --tags)
.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	CGO_ENABLED=0 go vet ./...
.PHONY:vet

build: vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/nct -ldflags "-s -w -X main.command=${COMMAND} -X main.version=${GIT_VER}"
.PHONY:build

install:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags "-s -w -X main.command=${COMMAND} -X main.version=${GIT_VER}"
.PHONY:install