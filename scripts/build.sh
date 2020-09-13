#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/linux/dgbr .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  -a -installsuffix cgo -o bin/macOS/dgbr .