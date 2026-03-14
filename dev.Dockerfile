FROM golang:1.25.6-alpine3.23 AS builder

ENV ENV=container
ENV CONFIG_PATH=/app/config
EXPOSE 8081

WORKDIR /app

COPY ./go.mod  ./
COPY ./go.sum  ./

RUN go mod download

COPY . .

ENTRYPOINT ["go", "tool", "air"]
