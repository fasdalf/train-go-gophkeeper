#!/bin/bash

go test ./... -bench=. -coverprofile cover.tmp.out && \
cat cover.tmp.out | grep -v ".pb.go:" | grep -v "/internal/client/" | grep -v "/train-go-gophkeeper/cmd/" > cover.out && \
go tool cover -html=cover.out -o cover.out.html && \
go tool cover -func cover.out > cover.log

# go test -v -covermode=count -coverprofile=coverage.out ./...