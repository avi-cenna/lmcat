#!/usr/bin/env just --justfile

set windows-shell := ["nu", "-c"]
set ignore-comments := true

fmt:
  go fmt .
#  golines . -w

lint:
  golangci-lint run

update:
  go get -u
  go mod tidy -v

build:
  go build

install:
  go build
  cp lmcat ~/bin/

run: build
  ./lmcat

stats: build
  ./lmcat --stats

demo:
  go build -ldflags "-s -w"
  ./lmcat --gcw
  hyperfine './lmcat --gcw'

demo-count: build
  cat ./common.go | ./lmcat --count
