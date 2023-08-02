docker/cloud/build:
	docker build -f docker/cloud/Dockerfile -t go-nitro-cloud .

docker/cloud/push:
	docker tag go-nitro-cloud:latest registry.digitalocean.com/magmo/go-nitro:latest
	docker push registry.digitalocean.com/magmo/go-nitro:latest

docker/local/build:
	docker build -f docker/local/Dockerfile -t go-nitro-local .

docker/local/start:
	docker remove go-nitro-local || true
	docker run -it -d --name go-nitro-local -p 3005:3005 -p 4005:4005 go-nitro-local

docker/build:
	docker build -f docker/Dockerfile -t go-nitro .

docker/start:
	docker remove go-nitro || true
	docker run -it -d --name go-nitro -p 4005:4005 -p 4006:4006 -p 4007:4007 go-nitro

docker/stop:
	docker stop go-nitro
	docker rm go-nitro

docker/restart: docker/stop docker/start

docker/attach:
	docker exec -it go-nitro bash

ui/build:
	yarn workspace nitro-gui build
