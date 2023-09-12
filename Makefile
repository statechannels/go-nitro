docker/cloud/build:
	docker build -f docker/cloud/Dockerfile -t go-nitro-cloud .

docker/cloud/start:
	docker remove go-nitro-cloud || true
	docker run -it -d --name go-nitro-cloud -p 3005:3005 -p 4005:4005 -p 5005:5005 go-nitro-cloud

docker/cloud/push:
	docker tag go-nitro-cloud:latest registry.digitalocean.com/magmo/go-nitro:latest
	docker push registry.digitalocean.com/magmo/go-nitro:latest

docker/local/build:
	docker build -f docker/local/Dockerfile -t go-nitro-local .

docker/local/start:
	docker remove go-nitro-local || true
	docker run -it -d --name go-nitro-local -p 3005:3005 -p 4005:4005 -p 5005:5005 go-nitro-local

docker/network/build:
	docker build -f docker/Dockerfile -t go-nitro .

docker/network/start:
	docker remove go-nitro || true
	docker run -it -d --name go-nitro -p 4005:4005 -p 4006:4006 -p 4007:4007 go-nitro

docker/network/stop:
	docker stop go-nitro
	docker rm go-nitro

docker/network/restart: docker/network/stop docker/network/start

docker/network/attach:
	docker exec -it go-nitro bash

ui/build:
	cd packages/nitro-rpc-client && bun run prepack
	cd packages/nitro-gui && bun run build
