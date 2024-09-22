#!/usr/bin/env just --justfile

set windows-shell := ["nu", "-c"]

fmt:
 go fmt .

update:
  go get -u
  go mod tidy -v

install:
  go build
  cp lmcat.exe ~/bin/