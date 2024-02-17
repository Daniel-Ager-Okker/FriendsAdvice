package services

import (
	model "FriendsAdvice/internal/data-model"
	"encoding/json"
	"time"
)

type IStorageManager interface {
	PutDataWithExpires(data *model.Data, whenDelete time.Time) (bool, error)
	PutData(data *model.Data) (bool, error)
	GetData(dataID uint64) *model.Data
	IsReady() bool
	Terminate() (bool, error)
}

type Controller struct {
	storageManager IStorageManager
}

func (c *Controller) IsStorageReady() bool {
	return c.storageManager.IsReady()
}

func (c *Controller) PutObject(key uint64, value []byte) (bool, error) {
	data := &model.Data{}
	err := json.Unmarshal(value, data)
	if err != nil {
		return false, err
	}

	return c.storageManager.PutData(data)
}

func (c *Controller) PutObjectWithExpires(key uint64, value []byte, lifetime time.Duration) (bool, error) {
	data := &model.Data{}
	err := json.Unmarshal(value, data)
	if err != nil {
		return false, err
	}

	wheneNeedToDelete := time.Now().Add(lifetime)
	return c.storageManager.PutDataWithExpires(data, wheneNeedToDelete)
}

func (c *Controller) GetObject(key uint64) ([]byte, bool) {
	data := c.storageManager.GetData(key)
	if data == nil {
		return []byte{}, false
	}

	obj, err := json.Marshal(data)
	if err != nil {
		return []byte{}, false
	}
	return obj, true
}

func InitController(storeManager IStorageManager) *Controller {
	controller := Controller{storageManager: storeManager}
	return &controller
}
