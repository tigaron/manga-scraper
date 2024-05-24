#!/usr/bin/env bash

# Generate prisma client
go generate ./...

# Build the app
go build -tags netgo -ldflags '-s -w' -o app cmd/manga-scraper/main.go
