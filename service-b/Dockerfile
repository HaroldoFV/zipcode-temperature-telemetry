FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init -g internal/server/main.go -o ./docs

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service-b ./internal/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/service-b .

COPY --from=builder /app/.env .

COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./service-b"]