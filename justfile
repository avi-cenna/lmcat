#!/usr/bin/env just --justfile

set windows-shell := ["nu", "-c"]

fmt:
  go fmt .
  golines . -w

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