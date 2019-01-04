# Build image
FROM golang:1.11 AS builder

WORKDIR /slack-nats
COPY . ./
RUN go mod download
RUN go test ./...
RUN CGO_ENABLED=0 go build

# Run image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /slack-nats/slack-nats .
ENTRYPOINT "/slack-nats"