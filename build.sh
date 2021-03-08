#!/bin/sh

# Delete app directory 
#rm -f ./bin

# Build app
go build -o ./bin/app

# run app
./bin/app