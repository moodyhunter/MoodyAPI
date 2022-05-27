#!/bin/bash

go generate
go get -u ./...
go mod tidy
