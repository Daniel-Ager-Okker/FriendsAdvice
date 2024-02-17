package services

import (
	"FriendsAdvice/internal/database/postgresql"
	"testing"
)

func TestIsReady(t *testing.T) {
	storageManager, _ := postgresql.InitManager(createConnectionDTO())
	controller := InitController(storageManager)
	isReady := controller.IsStorageReady()
	if !isReady {
		t.Errorf("Error while initializing inner storage manager!")
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
