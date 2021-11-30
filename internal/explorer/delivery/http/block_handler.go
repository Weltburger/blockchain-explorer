package http

import (
	"explorer/internal/explorer"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type BlockHandler struct {
	blockUseCase explorer.BlockUseCase
}

func NewBlockHandler(blockUseCase explorer.BlockUseCase) *BlockHandler {
	return &BlockHandler{
		blockUseCase: blockUseCase,
	}
}

func (h *BlockHandler) GetBlock(c echo.Context) error {
	blk := c.Param("block")

	block, err := h.blockUseCase.GetBlock(c.Request().Context(), blk)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, block)
}

func (h *BlockHandler) GetBlocks(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 1
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	blocks, err := h.blockUseCase.GetBlocks(c.Request().Context(), offset, limit)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, blocks)
}
