DOCKER_COMPOSE = docker-compose
DB_CONNECTION = "gabriel:password@tcp(localhost:3306)/snapfi?charset=utf8mb4&parseTime=True&loc=Local"

run:
	$(DOCKER_COMPOSE) up -d
	go run cmd/snapfi/main.go

stop:
	$(DOCKER_COMPOSE) down

test:
	go clean -testcache
	go test -v -p 1 -cover -failfast ./... -coverprofile=coverage.out
	go tool cover -func coverage.out | awk 'END{print sprintf("coverage: %s", $$3)}'

test-cover: test
	go tool cover -html=coverage.out 

mig-up:
	goose -dir ./internal/migrations mysql $(DB_CONNECTION) up

mig-down:
	goose -dir ./internal/migrations mysql $(DB_CONNECTION) down

mocks:
	rm -rf ./internal/mocks

	mockgen -source=./internal/database/user/user.go -destination=./internal/mocks/user.go -package=mocks -mock_names=Database=MockUserDatabase
	mockgen -source=./internal/database/transaction/transaction.go -destination=./internal/mocks/transaction.go -package=mocks -mock_names=Database=MockTransactionDatabase

	mockgen -source=./internal/app/user/user.go -destination=./internal/mocks/user_app.go -package=mocks -mock_names=App=MockUserApp
	mockgen -source=./internal/app/transaction/transaction.go -destination=./internal/mocks/transaction_app.go -package=mocks -mock_names=App=MockTransactionApp
