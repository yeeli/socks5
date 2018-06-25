#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o ./build/airship .
go build -o ./build/airship_mac .
