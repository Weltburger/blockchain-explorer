package http

import (
	"database/sql"
	"explorer/internal/explorer/repository/clickhouse"
	"explorer/internal/explorer/usecase"
	"github.com/labstack/echo/v4"
)

func RegisterEndpoints(router *echo.Echo, db *sql.DB) {
	endpoints := router.Group("/v1")
	{
		endpoints.GET("/blocks", NewBlockHandler(usecase.NewBlockUseCase(clickhouse.NewBlockRepository(db))).GetBlocks)
		endpoints.GET("/block/:block", NewBlockHandler(usecase.NewBlockUseCase(clickhouse.NewBlockRepository(db))).GetBlock)

		endpoints.GET("/transactions", NewTransHandler(clickhouse.NewTransRepository(db)).GetTransactions)
	}
}
