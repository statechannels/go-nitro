docker/nitro/build:
	docker build -f docker/cloud/Dockerfile -t go-nitro .

docker/nitro/start:
	docker remove go-nitro-cloud || true
	docker run -it -d --name go-nitro \
	-p 3005:3005 -p 4005:4005 -p 5005:5005 \
	-e NITRO_CONFIG_PATH="./nitro_config/iris.toml" \
	-v ./docker/nitro:/app/nitro_config \
	go-nitro

docker/nitro/push:
	docker tag go-nitro:latest registry.digitalocean.com/magmo/go-nitro:latest
	docker push registry.digitalocean.com/magmo/go-nitro:latest

docker/paymentproxy/build:
	docker build -f docker/paymentproxy/Dockerfile -t nitro-payment-proxy .

docker/paymentproxy/push:
	docker tag nitro-payment-proxy:latest registry.digitalocean.com/magmo/nitro-payment-proxy:latest
	docker push registry.digitalocean.com/magmo/nitro-payment-proxy:latest
	
docker/paymentproxy/start:
	docker remove payment-proxy || true
	docker run -it -d --name payment-proxy -p 5511:5511 -e PROXY_PORT=5511 payment-proxy

ui/build:
	yarn workspace nitro-gui build
