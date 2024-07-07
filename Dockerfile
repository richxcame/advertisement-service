FROM golang:1.22 AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download && go mod verify
COPY . .
COPY config.example.yml config.yml
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./bin/ads ./cmd/server

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/bin/ads /ads
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/config.yml /config.yml
CMD ["/ads"]