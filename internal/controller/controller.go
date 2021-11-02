package controller

import (
	"explorer/internal/storage"
)

type Controller struct {
	DB              *storage.Database
	blockController *BlockController
}

func (controller *Controller) BlockController() *BlockController {
	if controller.blockController != nil {
		return controller.blockController
	}

	controller.blockController = &BlockController{controller: controller}

	return controller.blockController
}

func New() *Controller {
	db := storage.GetDB()
	db.Tx, _ = db.DB.Begin()
	db.Stmt = db.BlockStorage().Prepare()
	return &Controller{DB: db}
}
