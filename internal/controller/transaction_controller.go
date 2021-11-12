package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type TransactionController struct {
	controller *Controller
}

func (transactionController *TransactionController) GetTransactions(c echo.Context) error {
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

	transactions, err := transactionController.controller.DB.TransactionStorage().GetTransactions(offset, limit, blk, hash, acc)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, transactions)
}
