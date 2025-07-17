FROM golang:bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o=./bin/app ./cmd

RUN apt-get update && apt-get install -y --no-install-recommends build-essential ca-certificates

RUN cd ./internal/repo/fasttext/fastText-0.9.2 && \
    make && mv ./fasttext ../../../../bin

RUN cd ./internal/repo/fasttext && \
    ../../../bin/fasttext supervised \
    -input ./labels.txt \
    -output ../../../bin/model \
    -epoch 100 \
    -dim 100 \
    -lr 0.10

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=builder /app/bin /app/bin

ENTRYPOINT ["./bin/app"]