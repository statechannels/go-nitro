FROM golang:1.20-bullseye

WORKDIR /go-nitro
COPY . .
RUN go build ./...
RUN chmod +x ./docker/entrypoint.sh

EXPOSE 4005 4006 4007

ENTRYPOINT ["./docker/entrypoint.sh"]
