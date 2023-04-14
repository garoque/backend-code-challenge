package transaction

import (
	"math"
	"net/http"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/labstack/echo/v4"
)

func Register(router *echo.Group, app *app.Container) {
	h := &handler{app}

	router.POST("", h.create)
	router.PUT("/increase-balance", h.increaseBalance)
}

type handler struct {
	app *app.Container
}

// Create transaction godoc
// @Summary Create transaction
// @Description Create transaction
// @Tags transaction
// @Accept json
// @Produce json
// @Param request body dto.CreateTransaction true "transaction request"
// @Success 201 {object} entity.Transaction
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /transaction [post]
func (h *handler) create(c echo.Context) error {
	var transaction dto.CreateTransaction
	if err := c.Bind(&transaction); err != nil {
		return echo.ErrInternalServerError
	}

	if err := c.Validate(&transaction); err != nil {
		return echo.ErrBadRequest
	}

	if math.Signbit(transaction.Amount) {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "The provided value is zero or negative")
	}

	balance, err := h.app.Transaction.Create(c.Request().Context(), entity.NewTransaction(transaction))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.Response{Data: balance})
}

// Increase balance user godoc
// @Summary Increase balance user
// @Description Increase balance user
// @Tags transaction
// @Accept json
// @Produce json
// @Param request body dto.IncreaseBalanceUser true "increase balance request"
// @Success 200 {object} float64
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /transaction/increase-balance [put]
func (h *handler) increaseBalance(c echo.Context) error {
	var transaction dto.IncreaseBalanceUser
	if err := c.Bind(&transaction); err != nil {
		return echo.ErrInternalServerError
	}

	if err := c.Validate(&transaction); err != nil {
		return echo.ErrBadRequest
	}

	if math.Signbit(transaction.Value) {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "The provided value is zero or negative")
	}

	balance, err := h.app.Transaction.IncreaseBalanceUser(c.Request().Context(), entity.NewIncreaseBalanceUser(transaction))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.Response{Data: balance})
}
