package postgresql

import (
	model "FriendsAdvice/internal/data-model"
	"FriendsAdvice/internal/services"
	"testing"
)

// Unfortunately, you may test this module only while have deployed database in docker container
// In this case we have some pseudo data for check postgresql package functionality

// Constants for connection to DB
// Be sure that user, password, port and dbname are the same as in docker-compose file
// host - depends on your own machine (TODO: make logic compile definitions)
const (
	host     string = "172.25.32.1"
	user     string = "postgres"
	password string = "postgres"
	port     string = "5432"
	dbname   string = "postgres"

	DANIEL_KEY uint = 1
	MARIA_KEY  uint = 2
	ZIMA_KEY   uint = 3
)

func makeTestConnectionDTO() *ConnectionDTO {
	return &ConnectionDTO{host, user, password, port, dbname}
}

func TestInitManager(t *testing.T) {
	// Check creation
	storageManager, err := InitManager(makeTestConnectionDTO())
	if err != nil {
		t.Fatal(err)
	}

	// Check connection to the database
	readyStatus := storageManager.IsReady()
	if !readyStatus {
		t.Errorf("Ready status for manager is wrong. Have: %v, want: %v.", readyStatus, true)
	}

	// Check inner default data
	checkDefaultData(t, storageManager)
}

func checkDefaultData(t *testing.T, storageManager services.IStorageManager) {
	// Daniel
	var dataDaniel *model.Data = storageManager.GetData(DANIEL_KEY)
	if !checkData(dataDaniel, "AMTEK", model.Physics, model.Test, 5) {
		t.Errorf("Daniel's default data is wrong!")
	}

	// Maria
	var dataMaria *model.Data = storageManager.GetData(MARIA_KEY)
	if !checkData(dataMaria, "School 26", model.English, model.Work, 4) {
		t.Errorf("Maria's default data is wrong!")
	}

	// Zima
	var dataZima *model.Data = storageManager.GetData(ZIMA_KEY)
	if !checkData(dataZima, "School 17", model.Mathematics, model.Quarter, 2) {
		t.Errorf("Maria's default data is wrong!")
	}
}

func checkData(data *model.Data, e string, s model.SubjectType, c model.KnowlegdeTestType, g model.GradeType) bool {
	cond1 := data.Establishment == e
	cond2 := data.Subject == s
	cond3 := data.KnowlegdeTest == c
	cond4 := data.Grade == g
	return cond1 && cond2 && cond3 && cond4
}
