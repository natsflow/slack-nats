# Build image
FROM golang:1.11 AS builder

ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

WORKDIR $GOPATH/src/github.com/natsflow/slack-nats

COPY . ./

RUN dep ensure -vendor-only
RUN go test ./...
RUN CGO_ENABLED=0 go install

# Run image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/slack-nats .
ENTRYPOINT "/slack-nats"