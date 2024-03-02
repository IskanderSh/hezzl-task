HOST = localhost

run:
	go run cmd/main.go --config=./config/local.yaml

tidy:
	go mod tidy

migrations-up:
	goose -dir "./migrations" postgres "host=${HOST} port=5432 user=postgres password=password" up

app-up:
	docker build -t application .
	docker run --rm \
	--name application \
	-p 1111:1111 \
	-d application

docker-up:
	docker-compose up -d