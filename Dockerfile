FROM golang:1.25-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]