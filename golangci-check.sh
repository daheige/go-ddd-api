#!/bin/bash
go mod tidy
golangci-lint run ./... > golangci.log
