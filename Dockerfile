FROM golang:1.27.0-alpine3.23 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o main main.go

FROM alpine:3.23

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 65532 appuser

ENV ENV=container
ENV CONFIG_PATH=/app/config
EXPOSE 8081

WORKDIR /app

COPY --from=builder /build/main .
COPY ./config /app/config

RUN chown -R appuser:appuser /app
USER appuser

ENTRYPOINT ["./main"]
