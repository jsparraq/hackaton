#!/bin/sh

# Delete app directory 
rm -rf ./bin

# Build app
go build -o ./bin/app

# run app
./bin/app