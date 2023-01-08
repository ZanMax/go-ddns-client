#!/bin/bash
# Build the project for different platforms

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/ddns-linux ddns.go

# Build for Linux arch ARM
GOOS=linux GOARCH=arm64 go build -o bin/ddns-linux-arm64 ddns.go