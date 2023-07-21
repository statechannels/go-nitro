docker/build:
	docker build -f docker/Dockerfile -t go-nitro .

docker/start:
	docker run -it -d --name go-nitro -p 4005:4005 -p 4006:4006 -p 4007:4007 go-nitro

docker/stop:
	docker stop go-nitro
	docker rm go-nitro

docker/restart: docker/stop docker/start

docker/attach:
	docker exec -it go-nitro bash

ui/build:
	yarn workspace nitro-gui build