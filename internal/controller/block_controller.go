package controller

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type BlockController struct {
	controller *Controller
}

func (blockController *BlockController) GetBlock(c echo.Context) error {
	blk := c.Param("block")

	ctx := context.Background()
	block, err := blockController.controller.DB.BlockStorage().GetBlock(ctx, blk)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	fmt.Println(block)

	return c.JSON(http.StatusOK, block)
}

func (blockController *BlockController) GetBlocks(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 1
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	ctx := context.Background()
	blocks, err := blockController.controller.DB.BlockStorage().GetBlocks(ctx, offset, limit)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, blocks)
}
