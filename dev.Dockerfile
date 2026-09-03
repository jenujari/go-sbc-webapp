FROM golang:1.27.0-alpine3.23

RUN apk add --no-cache git ca-certificates

ENV ENV=container
ENV CONFIG_PATH=/app/config
EXPOSE 8081

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENTRYPOINT ["go", "tool", "air"]
