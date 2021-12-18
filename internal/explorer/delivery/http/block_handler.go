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

// @Summary GetBlock
// @Security ApiKeyAuth
// @Tags blocks
// @Description Get block by hash or id
// @ID get-block
// @Produce  json
// @Success 200 {object} models.Block
// @Failure 404 {string} string "error"
// @Param block path string true "Block id or hash"
// @Router /v1/block/{block} [get]
func (h *BlockHandler) GetBlock(c echo.Context) error {
	blk := c.Param("block")

	block, err := h.blockUseCase.GetBlock(c.Request().Context(), blk)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, block)
}

// @Summary GetBlocks
// @Security ApiKeyAuth
// @Tags blocks
// @Description Get blocks with limit and offset
// @ID get-blocks
// @Produce  json
// @Success 200 {array} models.Block
// @Failure 404 {string} string "error"
// @Param limit  query int false "the amount of blocks you want to get"
// @Param offset query int false "offset from the beginning of the data in database"
// @Router /v1/blocks [get]
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
