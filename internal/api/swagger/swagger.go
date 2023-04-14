package swagger

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Register(router *echo.Group) {
	router.GET("/*", echoSwagger.WrapHandler)
}
