package services

import (
	model "FriendsAdvice/internal/data-model"
	"FriendsAdvice/internal/database/postgresql"
	"encoding/json"
	"testing"
	"time"
)

func TestIsReady(t *testing.T) {
	storageManager, _ := postgresql.InitManager(createConnectionDTO())
	controller := InitController(storageManager)
	isReady := controller.IsStorageReady()
	if !isReady {
		t.Errorf("Error while initializing inner storage manager!")
	}
}

func TestGetData(t *testing.T) {
	storageManager, _ := postgresql.InitManager(createConnectionDTO())
	controller := InitController(storageManager)

	jData, got := controller.GetObject(1)
	if !got {
		t.Errorf("Error while getting exisitng data!")
	}

	data := &model.Data{}
	err := json.Unmarshal(jData, data)
	if err != nil {
		t.Errorf("Error while unmarshaling exisitng data!")
	}

	if !check(data, "Daniel", "AMTEK", model.Eleven, model.LetterA) {
		t.Errorf("Got invalid data!")
	}
}

func TestPutData(t *testing.T) {
	storageManager, _ := postgresql.InitManager(createConnectionDTO())
	controller := InitController(storageManager)

	someData := model.Data{ID: 557,
		Pupil:         "Andrey",
		Establishment: "SHGD",
		Class:         model.Seven,
		Letter:        model.LetterV}

	someObj, err := json.Marshal(someData)
	if err != nil {
		t.Errorf("Error while marshaling some data!")
	}

	put, err := controller.PutObject(someData.ID, someObj)
	if err != nil || !put {
		t.Errorf("Error while putting some data!")
	}

	_, got := controller.GetObject(someData.ID)
	if !got {
		t.Errorf("Error while getting newly created data!")
	}
}

func TestPutDataWithExpires(t *testing.T) {
	storageManager, _ := postgresql.InitManager(createConnectionDTO())
	controller := InitController(storageManager)

	someData := model.Data{ID: 999,
		Pupil:         "Nelly",
		Establishment: "School 87",
		Class:         model.Eleven,
		Letter:        model.LetterA}

	someObj, err := json.Marshal(someData)
	if err != nil {
		t.Errorf("Error while marshaling some data!")
	}

	put, err := controller.PutObjectWithExpires(someData.ID, someObj, time.Nanosecond)
	if err != nil || !put {
		t.Errorf("Error while putting with expires some data!")
	}

	time.Sleep(time.Microsecond * 3)

	_, got := controller.GetObject(someData.ID)
	if got {
		t.Errorf("Error while getting newly created data with expires!")
	}
}

func createConnectionDTO() *postgresql.ConnectionDTO {
	return &postgresql.ConnectionDTO{
		HostName: "172.25.32.1",
		User:     "postgres",
		Password: "postgres",
		Port:     "5432",
		DBName:   "postgres"}
}

func check(data *model.Data, p, e string, c model.ClassType, l model.LetterType) bool {
	cond1 := data.Pupil == p
	cond2 := data.Establishment == e
	cond3 := data.Class == c
	cond4 := data.Letter == l
	return cond1 && cond2 && cond3 && cond4
}
