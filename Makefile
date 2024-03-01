HOST = localhost

run:
	go run cmd/main.go --config=./config/local.yaml

tidy:
	go mod tidy

postgres-up:
	docker run \
	-p 5432:5432 \
    --name postgres \
	-e POSTGRES_PASSWORD=password \
	-d postgres:latest

postgres-down:
	docker stop postgres
	docker rm postgres

postgres-connect:
	docker exec -it postgres /bin/bash

migrations-up:
	goose -dir "./migrations" postgres "host=${HOST} port=5432 user=postgres password=password" up

app-up:
	docker build -t application .
	docker run \
	--name application \
	-p 1111:1111 \
	-d application

docker-up:
	docker-compose up -d