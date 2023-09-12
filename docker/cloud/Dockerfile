FROM golang:1.21-bullseye AS builder

# Install Node.js v18.x
RUN apt-get update && apt-get install -y curl
RUN curl -sL https://deb.nodesource.com/setup_18.x | bash -
RUN apt-get install -y nodejs

# Install Yarn v3.5.1
RUN corepack enable
RUN yarn set version 3.5.1
RUN yarn --version

# Copy files into image
WORKDIR /app
COPY . .
COPY ./docker/cloud ./nitro_config

# Build the binary
RUN yarn workspace nitro-gui install
RUN yarn workspace nitro-gui build
RUN go build -v -o nitro -tags embed_ui .

# Reduce image size
FROM debian:bullseye-slim
RUN apt-get update
RUN apt-get install -y ca-certificates
RUN rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/nitro .
COPY --from=builder /app/nitro_config ./nitro_config

EXPOSE 3005 4005 5005

CMD ./nitro --config $CONFIG_PATH
