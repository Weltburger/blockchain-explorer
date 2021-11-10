package controller

import (
	"explorer/internal/storage"
)

type Controller struct {
	DB                    *storage.Database
	blockController       *BlockController
	transactionController *TransactionController
}

func (controller *Controller) BlockController() *BlockController {
	if controller.blockController != nil {
		return controller.blockController
	}

	controller.blockController = &BlockController{controller: controller}

	return controller.blockController
}

func (controller *Controller) TransactionController() *TransactionController {
	if controller.transactionController != nil {
		return controller.transactionController
	}

	controller.transactionController = &TransactionController{controller: controller}

	return controller.transactionController
}

func New() *Controller {
	db := storage.GetDB()
	db.BlockTx, _ = db.DB.Begin()
	db.BlockStmt = db.BlockStorage().Prepare()
	db.TransactionTx, _ = db.DB.Begin()
	db.TransactionStmt = db.TransactionStorage().Prepare()
	return &Controller{DB: db}
}
