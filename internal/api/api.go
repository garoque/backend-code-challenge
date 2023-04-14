package api

import (
	"github.com/garoque/backend-code-challenge-snapfi/internal/api/swagger"
	"github.com/garoque/backend-code-challenge-snapfi/internal/api/transaction"
	"github.com/garoque/backend-code-challenge-snapfi/internal/api/user"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app"
	"github.com/labstack/echo/v4"
)

func Register(router *echo.Group, app *app.Container) {
	user.Register(router.Group("/user"), app)
	transaction.Register(router.Group("/transaction"), app)
	swagger.Register(router.Group("/swagger"))
}
