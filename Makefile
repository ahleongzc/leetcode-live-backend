.DEFAULT: run
.PHONY: run build gen vet fmt count train fasttext compile

run: setUpDev compile gen build
	@if [ ! -f ./bin/model.bin ]; then \
        echo "model.bin not found. Running training..."; \
        make train; \
    fi
	@./bin/app

fasttext:
	@if [ ! -f ./bin/fasttext ]; then \
        echo "fasttext binary not found. running make from source..."; \
		cd ./internal/repo/fasttext/fastText-0.9.2 && make; \
		cp ./fasttext ../../../../bin/; \
    fi

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
	@cloc . --exclude-dir=scripts,.venv,fastText-0.9.2

train:
	@cd scripts && python3 clean.py
	@cd ./internal/repo/fasttext && \
		../../../bin/fasttext supervised \
		-input ./labels.txt \
		-output ../../../bin/model \
		-epoch 100 \
		-dim 100 \
		-lr 0.10

testModel:
	@./internal/repo/fasttext/fasttext test \
		./bin/model.bin \
		./scripts/test.txt

setUpDev:
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

apply:
	@cd infra && terraform apply

destroy:
	@cd infra && terraform destroy

compile:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./pb/*.proto