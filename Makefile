.DEFAULT: run
.PHONY: run build gen vet fmt count train

run: setUpInfra build
	@if [ ! -f ./bin/model.bin ]; then \
        echo "model.bin not found. Running training..."; \
        make train; \
    fi
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
	@cloc . --exclude-dir=scripts,.venv

train:
	@cd scripts && python3 clean.py
	@./internal/repo/fasttext/fasttext supervised \
		-input ./scripts/labels.txt \
		-output ./bin/model \
		-epoch 100 \
		-dim 100 \
		-lr 0.10

testModel:
	@./internal/repo/fasttext/fasttext test \
		./bin/model.bin \
		./scripts/test.txt

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