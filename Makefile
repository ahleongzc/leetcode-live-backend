.DEFAULT: run
.PHONY: run build gen vet fmt count

run: build
	@air

build: tidy fmt gen
	@go build -o=./bin/app ./cmd

gen:
	@cd internal/wire && wire

tidy:
	@go mod tidy

vet:
	@go vet ./...

fmt:
	@go fmt ./...

count:
	@cloc .