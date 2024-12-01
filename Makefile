.DEFAULT_GOAL := run

fmt: 
	@go fmt ./...

build: fmt
	@mkdir -p bin
	@go build -o bin/task main.go

run: build
	@./bin/task
