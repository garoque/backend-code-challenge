package user

import (
	"net/http"
	"strings"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/labstack/echo/v4"
)

var (
	ERROR_BAD_REQUEST = echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	ERROR_INTERNAL    = echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
)

func Register(router *echo.Group, app *app.Container) {
	h := &handler{app}

	router.GET("", h.readAll)
	router.GET("/:id", h.readOne)
	router.POST("", h.create)
}

type handler struct {
	app *app.Container
}

func (h *handler) create(c echo.Context) error {
	var request dto.CreateUser
	if err := c.Bind(&request); err != nil {
		return echo.ErrInternalServerError
	}

	if err := c.Validate(&request); err != nil {
		return echo.ErrBadRequest
	}

	err := h.app.User.Create(c.Request().Context(), entity.NewUser(request))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.Response{Data: nil})
}

func (h *handler) readOne(c echo.Context) error {
	userId := c.Param("id")
	if strings.Trim(userId, " ") == "" {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "The provided ID is empty")
	}
	user, err := h.app.User.ReadOneById(c.Request().Context(), userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.Response{Data: user})
}

func (h *handler) readAll(c echo.Context) error {
	users, err := h.app.User.ReadAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.Response{Data: users})
}
