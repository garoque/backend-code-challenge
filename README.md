docker-compose
go run cmd/snapfi/main.go
goose para as migrations

rodar migration: goose -dir ./internal/migrations mysql "gabriel:password@tcp(localhost:3306)/snapfi?charset=utf8mb4&parseTime=True&loc=Local" up


1° docker-compose
2° migrations
3° go run cmd/snapfi/main.go






OBS: VER PRA QUE SERVE OU REMOVER volumes DO DOCKER COMPOSE