# load env variables from env file
# https://unix.stackexchange.com/a/348432
include .env
export

run_dev:
	go run cmd/server.go

test:
	go test ./... -v

build: test
	go build -o server ./cmd

run: build
	./server
