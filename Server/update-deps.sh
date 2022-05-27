#!/bin/bash

go get -u ./...
go generate
go mod tidy
