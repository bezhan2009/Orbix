#!/bin/bash
go build -gcflags=all="-B" -ldflags="-s -w" orbix.go
./orbix.exe
