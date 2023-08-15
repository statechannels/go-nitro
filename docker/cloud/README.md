# Running a go-nitro bootnode in the cloud

The Dockerfile in this directory defines a container that can act as a `bootnode` to boostrap a nitro network. When other nodes initialize, they can connect to this `bootnode` and use it to help discover peers. The following steps document how to deploy this node to a Digital Ocean Droplet VM. The config file and steps below assume you want to run the `bootnode` against the Filecoin Calibration testnet, but can be altered to connect to other blockchain networks.

## Create image

Update the `docker/cloud/config.toml` file with the private keys you will use for the state channel address and chain address. You can run the following command to generate a fresh key pair:

```
go run ./cmd/generate-keypair
```

You can use [this faucet](https://faucet.calibration.fildev.network/funds.html) to fund the chain address with some Filecoin Calibration testnet FIL:

Build the image:

```
make docker/cloud/build
```

## Push image to docker registry

```
docker login -u <api_key> -p <api_key> registry.digitalocean.com
make docker/cloud/push
```

## Run container in the cloud

Connect to the Droplet's terminal via ssh:

```
ssh -i <ssh_private_key> root@<droplet_ip_address>
```

Run the following commands from the Droplet's terminal to start the container:

```
docker pull registry.digitalocean.com/magmo/go-nitro:latest
docker run -it -d -p 3005:3005 -p 4005:4005 registry.digitalocean.com/magmo/go-nitro:latest
docker logs -f <container_id>
```
