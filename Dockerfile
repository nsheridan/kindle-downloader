FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o kindle-downloader .

FROM ghcr.io/go-rod/rod

COPY --from=builder /app/kindle-downloader /

ENTRYPOINT [ "/kindle-downloader" ]
