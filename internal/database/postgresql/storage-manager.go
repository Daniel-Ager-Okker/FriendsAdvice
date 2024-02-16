package postgresql

import (
	model "FriendsAdvice/internal/data-model"
	"FriendsAdvice/internal/services"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// An implementation of IStorageManager interface for work with database
type StorageManager struct {
	storage     map[uint]*model.Data
	keysBuffer  []uint
	bufferLimit uint
	dataBase    *sql.DB
}

// StorageManager method for put data object into the RAM and update DB if need
func (pManager *StorageManager) PutData(data *model.Data) bool {
	// 1.Check if data with such ID already exists
	if _, exists := pManager.storage[data.ID]; exists {
		return false
	}

	// 2.Put into the RAM and update keysBuffer
	pManager.storage[data.ID] = data
	pManager.keysBuffer = append(pManager.keysBuffer, data.ID)

	// 3.Update databse if it is necessary
	if len(pManager.keysBuffer) == int(pManager.bufferLimit) {
		dbUpdated := pManager.updateDataBase()
		pManager.keysBuffer = make([]uint, 0)
		return dbUpdated
	}
	return true
}

// StorageManager method for get data object by its ID
func (pManager *StorageManager) GetData(dataID uint) *model.Data {
	// Check if there is data with such ID and then return it
	if data, exists := pManager.storage[dataID]; exists {
		return data
	} else {
		return nil
	}
}

// StorageManager method for answer the question of it is ready
func (pManager *StorageManager) IsReady() bool {
	if pManager.dataBase == nil {
		return false
	}

	err := pManager.dataBase.Ping()
	return err == nil
}

// Function for right termination while service down
func (pManager *StorageManager) Terminate() (bool, error) {
	// Need to move all objects from RAM to DB and close connection
	pManager.updateDataBase()
	err := pManager.dataBase.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

// Function for create StorageManager enitity in according with possible data in the database
func InitManager(connectionInfo *ConnectionDTO) (services.IStorageManager, error) {
	// 1.First step - try to open database
	pDataBase, err := sql.Open("postgres", getOpenConnectionStatement(connectionInfo))
	if err != nil {
		return nil, err
	}

	// 2.Second step - check if there is connection
	err = pDataBase.Ping()
	if err != nil {
		pDataBase.Close()
		return nil, err
	}

	// 3.Get data from DB to the RAM
	manager := StorageManager{storage: make(map[uint]*model.Data),
		keysBuffer:  make([]uint, 0),
		bufferLimit: 30,
		dataBase:    pDataBase}
	manager.getDataFromDB()

	return &manager, nil
}

func (pManager *StorageManager) getDataFromDB() {
	rows, err := pManager.dataBase.Query("SELECT * FROM Grades")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id            uint
		pupil         string
		establishment string
		subject       model.SubjectType
		knowlegdeTest model.KnowlegdeTestType
		grade         model.GradeType
	)

	for rows.Next() {
		err = rows.Scan(&id, &pupil, &establishment, &subject, &knowlegdeTest, &grade)
		if err != nil {
			return
		}
		data := model.Data{ID: id,
			Pupil:         pupil,
			Establishment: establishment,
			Subject:       subject,
			KnowlegdeTest: knowlegdeTest,
			Grade:         grade}
		pManager.storage[id] = &data
	}
}

func (pManager *StorageManager) updateDataBase() bool {
	//TODO
	return true
}

func getOpenConnectionStatement(pInfo *ConnectionDTO) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s port=%s dbname=%s sslmode=%s",
		pInfo.HostName, pInfo.User, pInfo.Password, pInfo.Port, pInfo.DBName, "disable")
}
