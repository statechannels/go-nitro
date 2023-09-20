# Running a go-nitro node in the cloud

The Dockerfile in this directory defines a container that can act as a `bootnode` to boostrap a nitro network or as a common node that connects to `bootnodes`, depending on the `bootpeers` defined in the `.toml` config file. The following steps document how to deploy one of theses node to a Digital Ocean Droplet VM. The config file and steps below assume you want to run the node against the Filecoin Calibration testnet, but can be altered to connect to other blockchain networks.

## Create image

Update the `docker/nitro/config.toml` file (or pass as env vars) with the private keys you will use for the state channel address and chain address. You can run the following command to generate a fresh key pair:

```
go run ./cmd/generate-keypair
```

You can use [this faucet](https://faucet.calibration.fildev.network/funds.html) to fund the chain address with some Filecoin Calibration testnet FIL:

Build the image:

```
make docker/nitro/build
```

## Push image to docker registry

```
docker login -u <api_key> -p <api_key> registry.digitalocean.com
make docker/nitro/push
```

## Create ssh key-pair

```
ssh-keygen -b 4096 -t rsa -f <output_file_path>
```

## Settings for Digital Ocean Droplet to host the node

- Select `Choose an image`
  - Select `Marketplace`
    - Select `Docker 23.0.6 on Ubuntu 22.04`
- Add ssh public key created in previous step

## Accessing nitro node logs from Droplet

From within the Digital Ocean console, follow these steps to `ssh` into the correct Droplet and stream the logs:

1. Project should be `go-nitro`
1. Select `Droplets` in left side-bar menu
1. Select name of the `Droplet` you want to access
1. Select `Access` in inner left side-bar menu
1. Click `Launch Droplet Console` blue button
1. Run `docker logs -f <container_name>`
