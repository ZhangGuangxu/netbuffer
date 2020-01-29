#!/bin/sh

gofmt -s -l -w *.go && go vet *.go && errcheck *.go && go build -i

exit 0
