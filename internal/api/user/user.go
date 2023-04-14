package user

import (
	"net/http"
	"strings"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/labstack/echo/v4"
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

// Create user godoc
// @Summary Create user
// @Description Create user
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.CreateUser true "user request"
// @Success 201
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /user [post]
func (h *handler) create(c echo.Context) error {
	var request dto.CreateUser
	if err := c.Bind(&request); err != nil {
		return echo.ErrInternalServerError
	}

	if err := c.Validate(&request); err != nil {
		return echo.ErrBadRequest
	}

	err := h.app.User.Create(c.Request().Context(), *entity.NewUser(request))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.Response{Data: nil})
}

// Read one user godoc
// @Summary Read one user
// @Description Read one user
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "user ID" Format(uuid)
// @Success 200 {object} entity.User
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /user/{id} [get]
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

// Read read all users godoc
// @Summary Read read all users
// @Description Read read all users
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} entity.User
// @Failure 500 {object} error
// @Router /user [get]
func (h *handler) readAll(c echo.Context) error {
	users, err := h.app.User.ReadAll(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.Response{Data: users})
}
