package http

import (
	"explorer/internal/explorer"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type TransMIHandler struct {
	transUseCase explorer.TransMIUseCase
}

func NewTransMIHandler(transMIUseCase explorer.TransMIUseCase) *TransMIHandler {
	return &TransMIHandler{
		transUseCase: transMIUseCase,
	}
}

// @Summary GetTransactionsMI
// @Security ApiKeyAuth
// @Tags transactions-main-info
// @Description Get transactions main info with limit and offset params (transactions related to hen)
// @ID get-transactions-main-info
// @Produce  json
// @Success 200 {array} models.TransactionMainInfo
// @Failure 404 {string} string "error"
// @Param limit  query int false "the amount of transactions you want to get"
// @Param offset query int false "offset from the beginning of the data in database"
// @Param block query string false "specifying a block"
// @Param account query string false "specifying an address of account (src or dst)"
// @Param hash query string false "specifying the transaction hash"
// @Router /v1/transactions-main-info [get]
func (h *TransMIHandler) GetTransactionsMI(c echo.Context) error {
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
		c.String(http.StatusNotFound, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, transactions)
}
