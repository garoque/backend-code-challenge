DOCKER_COMPOSE = docker-compose
DB_CONNECTION = "gabriel:password@tcp(localhost:3306)/snapfi?charset=utf8mb4&parseTime=True&loc=Local"

run:
	$(DOCKER_COMPOSE) up -d
	go run cmd/snapfi/main.go

stop:
	$(DOCKER_COMPOSE) down

mig-up:
	goose -dir ./internal/migrations mysql $(DB_CONNECTION) up

mig-down:
	goose -dir ./internal/migrations mysql $(DB_CONNECTION) down