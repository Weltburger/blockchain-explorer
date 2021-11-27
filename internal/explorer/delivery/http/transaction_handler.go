package http

import (
	"explorer/internal/explorer"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type TransHandler struct {
	transUseCase explorer.TransUseCase
}

func NewTransHandler(transUseCase explorer.TransUseCase) *TransHandler {
	return &TransHandler{
		transUseCase: transUseCase,
	}
}

func (h *TransHandler) GetTransactions(c echo.Context) error {
	blk := c.QueryParam("block")
	acc := c.QueryParam("account")
	hash := c.QueryParam("hash")

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 1
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil  || offset < 0{
		offset = 0
	}

	transactions, err := h.transUseCase.GetTransactions(c.Request().Context(), offset, limit, blk, hash, acc)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, transactions)
}
