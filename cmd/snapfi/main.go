package main

import (
	"log"

	_ "github.com/garoque/backend-code-challenge-snapfi/docs"
	"github.com/garoque/backend-code-challenge-snapfi/internal/api"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app/transaction"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app/user"
	"github.com/garoque/backend-code-challenge-snapfi/internal/database"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title           Snapfi Backend Code Challenge
// @version         1.0
// @termsOfService  http://swagger.io/terms/

// @contact.name   Gabriel F Roque
// @contact.url    https://github.com/garoque/backend-code-challenge-snapfi
// @contact.email  gfroque1@gmail.com

// @license.name  Snapfi Backend Code Challenge
// @license.url   https://github.com/garoque/backend-code-challenge-snapfi

// @host      localhost:1323
// @BasePath  /v1
func main() {
	e := echo.New()

	connDb, err := sqlx.Open("mysql", "gabriel:password@tcp(localhost:3306)/snapfi?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
	defer connDb.Close()

	e.Validator = validator.NewValidator()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := database.New(connDb)

	api.Register(e.Group("/v1"), &app.Container{
		User:        user.NewAppUser(db),
		Transaction: transaction.NewAppTransaction(db),
	})

	e.Logger.Fatal(e.Start(":1323"))
}
