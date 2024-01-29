#! /bin/sh

if [ -z "$1" ]; then
    echo "Please specify the file name"
    exit 1
fi

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $1