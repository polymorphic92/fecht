#!/usr/bin/env bash
path=$(git rev-parse --show-toplevel)

go build -ldflags "-s -w" -o fet $path
sudo mv fet /usr/local/bin
