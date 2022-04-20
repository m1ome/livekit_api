FROM golang:1.18-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
RUN go build -o livekit_api

FROM alpine:3.15

WORKDIR /app
COPY --from=builder /app/livekit_api ./livekit_api

ENTRYPOINT ["/app/livekit_api"]
