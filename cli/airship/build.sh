#!/bin/bash

go build -o ../../build/airship .
GOOS=linux GOARCH=amd64 go build -o ../../build/airship_linux .
