.DEFAULT: run
.PHONY: run build gen vet fmt count

run: setUpInfra build
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

setUpInfra:
	@if [ "$$(docker ps -q -f name=rabbitmq)" = "" ]; then \
		echo "Starting RabbitMQ container..."; \
		docker run -d \
			--name rabbitmq \
			-p 5672:5672 \
			-p 15672:15672 \
			-v rabbitmqData:/var/lib/rabbitmq \
			rabbitmq:4-management; \
	else \
		echo "RabbitMQ container already running."; \
	fi