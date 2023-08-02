FROM golang:1.20-bullseye AS builder
WORKDIR /app
COPY . .
RUN go build -o nitro .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/nitro .
COPY --from=builder /app/docker/cloud/config.toml .

EXPOSE 3005 4005

ENTRYPOINT ["./nitro --config ./config.toml"]
