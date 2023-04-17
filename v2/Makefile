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

# build: vet
# 	CGO_ENABLED=0 go build -o nct main.go
# .PHONY:build