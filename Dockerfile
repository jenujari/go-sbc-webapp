FROM golang:1.25.6-alpine3.23 AS builder

WORKDIR /build

COPY ./go.mod  ./
COPY ./go.sum  ./

RUN go mod download

COPY . .

RUN go build -o main main.go

FROM alpine:3.23

ENV ENV=container
ENV CONFIG_PATH=/app/config
EXPOSE 8081

WORKDIR /app

COPY --from=builder /build/main .
COPY ./config /app/config

ENTRYPOINT ["./main"]
