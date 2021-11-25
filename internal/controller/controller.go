package controller

import (
	"explorer/internal/storage"
)

type Controller struct {
	DB *storage.Database
}

func (controller *Controller) BlockController() *BlockController {
	return &BlockController{controller: controller}
}

func (controller *Controller) TransactionController() *TransactionController {
	return &TransactionController{controller: controller}
}

func New() *Controller {
	// db := storage.GetDB()
	db := &storage.Database{}

	return &Controller{DB: db}
}
