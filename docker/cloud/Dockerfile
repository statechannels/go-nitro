FROM golang:1.20-bullseye AS builder
WORKDIR /app
COPY . .
COPY ./docker/cloud/config.toml .
RUN go build -v -o nitro .

FROM debian:bullseye-slim
RUN apt-get update
RUN apt-get install -y ca-certificates
RUN rm -rf /var/lib/apt/lists/*
WORKDIR /app/
COPY --from=builder /app/nitro .
COPY --from=builder /app/config.toml .

EXPOSE 3005 4005

CMD ["./nitro", "--config", "./config.toml"]
