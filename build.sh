#!/bin/bash

echo "Building expressive"

go build -o $GOPATH/bin/expressive -v -a ./cli
