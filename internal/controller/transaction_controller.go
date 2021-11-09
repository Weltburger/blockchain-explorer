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
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 1
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}


	transactions, err := transactionController.controller.DB.TransactionStorage().GetTransactions(offset, limit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, transactions)
}

func (transactionController *TransactionController) GetTransactionsByBlock(c echo.Context) error {
	blk := c.Param("block")

	transactions, err := transactionController.controller.DB.TransactionStorage().GetTransactionsByBlock(blk)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, transactions)
}

func (transactionController *TransactionController) GetTransactionsByAddress(c echo.Context) error {
	addr := c.Param("address")

	transactions, err := transactionController.controller.DB.TransactionStorage().GetTransactionsByAddress(addr)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, transactions)
}

func (transactionController *TransactionController) GetTransactionsByHash(c echo.Context) error {
	hash := c.Param("hash")

	transactions, err := transactionController.controller.DB.TransactionStorage().GetTransactionsByHash(hash)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, transactions)
}
