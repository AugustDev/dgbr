#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/linux/dgbr .