package controller

import (
	"explorer/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type BlockController struct {
	controller *Controller
}

func (blockController *BlockController) GetBlock(c echo.Context) error {
	blk := c.Param("block")

	block := new(models.Block)

	resp, err := blockController.controller.DB.DB.Query(`
		SELECT * FROM block WHERE Hash = ?
	`, blk)
	if err != nil {
		return err
	}

	var tm time.Time
	var header, meatadata, ops string
	resp.Next()
	err = resp.Scan(&block.Protocol, &block.ChainID, &block.Hash, &tm, &header, &meatadata, &ops)
	if err != nil {
		return err
	}
	fmt.Println("aaaa", block, header, meatadata, ops)

	return c.String(http.StatusOK, "bloch has been received")
}
