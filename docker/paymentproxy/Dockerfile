FROM golang:1.21-bullseye AS builder
WORKDIR /app
COPY . .
COPY ./docker/local/config.toml .
RUN go build -o proxy ./cmd/start-payment-proxy

FROM debian:bullseye-slim
RUN apt-get update
RUN apt-get install -y ca-certificates
RUN rm -rf /var/lib/apt/lists/*
WORKDIR /app/
COPY --from=builder /app/proxy .


ENV PROXY_PORT=5511
ENV NITRO_ENDPOINT=host.docker.internal:4007/api/v1
ENV DESTINATION_URL=http://host.docker.internal:8088
ENV COST_PER_BYTE=1
ENV TLS_CERT_FILE=""
ENV TLS_KEY_FILE=""

EXPOSE $PROXY_PORT
CMD ./proxy --nitroendpoint $NITRO_ENDPOINT --proxyaddress 0.0.0.0:$PROXY_PORT --destinationurl $DESTINATION_URL --costperbyte $COST_PER_BYTE --tlscertfilepath $TLS_CERT_FILE --tlskeyfilepath $TLS_KEY_FILE 
