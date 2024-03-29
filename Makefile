HOST = localhost

run:
	go run cmd/main.go --config=.github/local/local.yaml

tidy:
	go mod tidy

migrations-up:
	goose -dir "./migrations" postgres "host=${HOST} port=5432 user=postgres password=password" up

app-up:
	docker build -t application -f Dockerfile.local
	docker run --rm \
	--name application \
	-p 1111:1111 \
	-d application

docker-up-local:
	docker-compose -f ./docker-compose-local.yml up -d

docker-up-prod:
	docker-compose -f ./docker-compose-prod.yml up -d