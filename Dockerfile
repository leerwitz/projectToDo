FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:3.20.1

WORKDIR /app
RUN apk add --no-cache bash

COPY --from=builder /app/main .
COPY wait-for-it.sh .

EXPOSE 8080

CMD ["./wait-for-it.sh", "db:5432", "--", "./main"]
