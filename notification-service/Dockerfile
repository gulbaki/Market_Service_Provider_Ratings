FROM golang:1.23-alpine AS builder 
WORKDIR /app

COPY . .
RUN go mod download
#It would be safer to do it in CI/CD.
RUN go test ./... -v
RUN go build -o notification-service ./cmd/notification-service

FROM alpine:3.17
WORKDIR /app

COPY --from=builder /app/notification-service .
EXPOSE 9191

ENTRYPOINT ["./notification-service"]
