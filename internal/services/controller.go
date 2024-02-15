package services

import (
	model "FriendsAdvice/internal/data-model"
	transport "FriendsAdvice/internal/transport"
	"time"
)

type IStorageManager interface {
	PutData(data *model.Data) bool
	GetData(dataID uint) *model.Data
	Terminate() (bool, error)
}

type Controller struct {
	storageManager IStorageManager
}

func (c *Controller) IsStorageReady() bool {
	// TODO
	return true
}

func (c *Controller) PutObject(key, value string) bool {
	// TODO
	return true
}

func (c *Controller) PutObjectWithExpires(key, value string, lifetime time.Duration) bool {
	// TODO
	return true
}

func (c *Controller) GetObject(key string) (string, bool) {
	// TODO
	return "ah", true
}

func InitController(storeManager IStorageManager) transport.IController {
	controller := Controller{storageManager: storeManager}
	return &controller
}
