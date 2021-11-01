package controller

import (
	"explorer/internal/storage"
	"explorer/models"
)

type Controller struct {
	DB              *storage.Database
	block           *models.Block
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
	return &Controller{DB: db}
}
