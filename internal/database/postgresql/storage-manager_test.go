package postgresql

import (
	model "FriendsAdvice/internal/data-model"
	"FriendsAdvice/internal/services"
	"testing"
	"time"
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
	if !checkData(dataDaniel, "Daniel", "AMTEK", model.Eleven, model.LetterA) {
		t.Errorf("Daniel's default data is wrong!")
	}

	// Maria
	var dataMaria *model.Data = storageManager.GetData(MARIA_KEY)
	if !checkData(dataMaria, "Maria", "School 26", model.Ten, model.LetterA) {
		t.Errorf("Maria's default data is wrong!")
	}

	// Zima
	var dataZima *model.Data = storageManager.GetData(ZIMA_KEY)
	if !checkData(dataZima, "Zima", "School 17", model.Ten, model.LetterV) {
		t.Errorf("Maria's default data is wrong!")
	}
}

func checkData(data *model.Data, p, e string, c model.ClassType, l model.LetterType) bool {
	cond1 := data.Pupil == p
	cond2 := data.Establishment == e
	cond3 := data.Class == c
	cond4 := data.Letter == l
	return cond1 && cond2 && cond3 && cond4
}

func TestPutData(t *testing.T) {
	storageManager, _ := InitManager(makeTestConnectionDTO())

	// some valid (means "new") UserData
	data := model.Data{ID: 4,
		Pupil:         "Antonio",
		Establishment: "AMTEK",
		Class:         model.Nine,
		Letter:        model.LetterA}
	added := storageManager.PutData(&data)
	if !added {
		t.Errorf("Newly created data put FAILED. Have: %v, want: %v.", added, true)
	}

	addedAgain := storageManager.PutData(&data)
	if addedAgain {
		t.Errorf("Existing data put FAILED. Have: %v, want: %v.", addedAgain, false)
	}
}

func TestExpires(t *testing.T) {
	storageManager, _ := InitManager(makeTestConnectionDTO())

	// some valid (means "new") UserData
	data := model.Data{ID: 7,
		Pupil:         "Karen",
		Establishment: "School 37",
		Class:         model.Five,
		Letter:        model.LetterB}
	added := storageManager.PutDataWithExpires(&data, time.Now())
	if !added {
		t.Errorf("Newly created data put FAILED. Have: %v, want: %v.", added, true)
	}

	gotData := storageManager.GetData(data.ID)
	if gotData != nil {
		t.Errorf("Expired data got, but should no!")
	}
}

func TestGetData(t *testing.T) {
	storageManager, _ := InitManager(makeTestConnectionDTO())

	// Try to get non-existen data
	dataNonExisten := storageManager.GetData(500)
	if dataNonExisten != nil {
		t.Errorf("Non-existen data get FAILED. Have: %v, want: nil.", dataNonExisten)
	}

	// Put some valid data
	data := model.Data{ID: 4,
		Pupil:         "Antonio",
		Establishment: "AMTEK",
		Class:         model.Nine,
		Letter:        model.LetterA}
	storageManager.PutData(&data)

	// Try to get it
	gotData := storageManager.GetData(data.ID)
	if gotData == nil {
		t.Errorf("Existing data get FAILED. Have: %v.", gotData)
	}
}

func TestTerminate(t *testing.T) {
	storageManager, _ := InitManager(makeTestConnectionDTO())

	// need to put some valid data and terminate manager
	// after that you should check database state manually
	storageManager.PutData(&model.Data{4, "Antonio", "AMTEK", model.Nine, model.LetterA})
	storageManager.PutData(&model.Data{5, "Olga", "School 26", model.Ten, model.LetterA})

	terminated, err := storageManager.Terminate()
	if err != nil {
		t.Fatal(err)
	}

	if !terminated {
		t.Errorf("Something wrong while terminating")
	}
}
