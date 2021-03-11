SERVER_IMAGE = reg.istry.cz/ondrejsika/example-pool-k8s-long-conn--server
CLIENT_IMAGE = reg.istry.cz/ondrejsika/example-pool-k8s-long-conn--client

docker:
	docker build --pull -t $(SERVER_IMAGE) ./server
	docker build --pull -t $(CLIENT_IMAGE) ./client
	docker push $(SERVER_IMAGE)
	docker push $(CLIENT_IMAGE)

