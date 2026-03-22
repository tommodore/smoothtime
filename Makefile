.PHONY: run test build
run:
go run .
build:
go build -o smoothtime .
docker:
docker build -t smoothtime .
