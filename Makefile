# load env variables from env file
# https://unix.stackexchange.com/a/348432
include .env
export

run_dev:
	go run main.go
