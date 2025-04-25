FROM golang:1.24.1 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

ENV CGO_ENABLED=0

RUN go build -o /app/app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

COPY ./config.json .

EXPOSE 8080

CMD ["./app"]