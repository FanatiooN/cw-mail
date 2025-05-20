FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init -g cmd/app/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mail-service ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/mail-service .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/entrypoint.sh .

RUN chmod +x entrypoint.sh

ENV GIN_MODE=release

EXPOSE 8080

CMD ["./entrypoint.sh"] 