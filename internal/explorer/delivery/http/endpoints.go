package http

import (
	"database/sql"
	"explorer/internal/explorer/repository/clickhouse"
	"github.com/labstack/echo/v4"
)

func RegisterEndpoints(router *echo.Echo, db *sql.DB) {
	endpoints := router.Group("/api/v1")
	{
		endpoints.GET("/blocks", NewBlockHandler(clickhouse.NewBlockRepository(db)).GetBlocks)
		endpoints.GET("/block/:block", NewBlockHandler(clickhouse.NewBlockRepository(db)).GetBlock)

		endpoints.GET("/transactions", NewTransHandler(clickhouse.NewTransRepository(db)).GetTransactions)
	}
}
