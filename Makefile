.DEFAULT: run
.PHONY: run build gen vet fmt count

run: build
	@air

build: fmt gen
	@go build -o=./bin/app ./cmd

gen:
	@cd internal/wire && wire

vet:
	@go vet ./...

fmt:
	@go fmt ./...

count:
	@cloc .