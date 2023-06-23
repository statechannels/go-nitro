FROM golang:1.20-bullseye

# Set up anvil
RUN curl -L https://foundry.paradigm.xyz | bash
ENV PATH="$PATH:/root/.foundry/bin"
RUN . /root/.bashrc
RUN foundryup

WORKDIR /go-nitro
COPY . .
RUN go build ./...
RUN chmod +x ./docker/entrypoint.sh

EXPOSE 4005 4006 4007

ENTRYPOINT ["./docker/entrypoint.sh"]
