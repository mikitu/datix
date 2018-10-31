#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./apps/api/main ./apps/api/main.go
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./apps/db/main ./apps/db/main.go