#!/bin/bash

# compile
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/linux/dgbr .
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  -a -installsuffix cgo -o bin/macOS/dgbr .

# gzip
gzip -c bin/linux/dgbr > release/dgbr.gz && mv release/dgbr.gz release/dgbr-linux-amd64.gz
gzip -c bin/macOS/dgbr > release/dgbr.gz && mv release/dgbr.gz release/dgbr-darwin-amd64.gz