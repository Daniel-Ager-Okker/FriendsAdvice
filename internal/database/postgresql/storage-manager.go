package postgresql

import (
	model "FriendsAdvice/internal/data-model"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	_ "github.com/lib/pq"
)

// An implementation of IStorageManager interface for work with database
type StorageManager struct {
	storage     map[uint64]*model.Data
	deleteInfo  treemap.Map
	keysBuffer  []uint64
	bufferLimit uint64
	dataBase    *sql.DB
}

// StorageManager method for put data with expires object into the RAM and update DB if need
func (pManager *StorageManager) PutDataWithExpires(data *model.Data, whenDelete time.Time) (bool, error) {
	dataPut, err := pManager.PutData(data)
	if !dataPut {
		return false, err
	}

	pManager.deleteInfo.Put(whenDelete, data.ID)

	// 4. Handle expired
	pManager.handleExpired()
	return true, nil
}

// StorageManager method for put data object into the RAM and update DB if need
func (pManager *StorageManager) PutData(data *model.Data) (bool, error) {
	// 1.Check if data with such ID already exists
	if _, exists := pManager.storage[data.ID]; exists {
		return false, errors.New("Already have such key in storage")
	}

	// 2.Put into the RAM and update keysBuffer
	pManager.storage[data.ID] = data
	pManager.keysBuffer = append(pManager.keysBuffer, data.ID)

	// 3.Update databse if it is necessary
	if len(pManager.keysBuffer) == int(pManager.bufferLimit) {
		dbUpdated, err := pManager.updateDataBase()
		pManager.keysBuffer = make([]uint64, 0)
		return dbUpdated, err
	}

	// 4. Handle expired
	pManager.handleExpired()
	return true, nil
}

// StorageManager method for get data object by its ID
func (pManager *StorageManager) GetData(dataID uint64) *model.Data {
	// Handle expired
	pManager.handleExpired()

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
	updated, err := pManager.updateDataBase()
	pManager.handleExpired()
	if !updated {
		return false, err
	}
	err = pManager.dataBase.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

// Function for delete every expired data entity
func (pManager *StorageManager) handleExpired() {
	var now time.Time = time.Now()
	keys := pManager.deleteInfo.Keys()
	lowerBoundIdx := lowerBound(keys, 0, len(keys), now)
	for i := 0; i < lowerBoundIdx; i++ {
		idForDelete, _ := pManager.deleteInfo.Get(keys[i].(time.Time))
		pManager.removeData(idForDelete.(uint64))
	}
}

func (pManager *StorageManager) removeData(ID uint64) {
	// Remove from RAM
	delete(pManager.storage, ID)

	// Remove from DB
	pManager.dataBase.Query("DELETE FROM Pupils WHERE id = $1", ID)
}

// Function for create StorageManager enitity in according with possible data in the database
func InitManager(connectionInfo *ConnectionDTO) (*StorageManager, error) {
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
	manager := StorageManager{storage: make(map[uint64]*model.Data),
		deleteInfo:  *treemap.NewWith(TimesComparator),
		keysBuffer:  make([]uint64, 0),
		bufferLimit: 30,
		dataBase:    pDataBase}
	manager.getDataFromDB()

	return &manager, nil
}

// Function for get data from db to RAM storage
func (pManager *StorageManager) getDataFromDB() {
	rows, err := pManager.dataBase.Query("SELECT * FROM Pupils")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var (
		id            uint64
		pupil         string
		establishment string
		classNum      model.ClassType
		letter        model.LetterType
	)

	for rows.Next() {
		err = rows.Scan(&id, &pupil, &establishment, &classNum, &letter)
		if err != nil {
			return
		}
		data := model.Data{ID: id,
			Pupil:         pupil,
			Establishment: establishment,
			Class:         classNum,
			Letter:        letter}
		pManager.storage[id] = &data
	}
}

// Function for move data from database to the RAM storage
func (pManager *StorageManager) updateDataBase() (bool, error) {
	for key, data := range pManager.storage {
		rows, err := pManager.dataBase.Query("SELECT * FROM Pupils WHERE id = $1", key)
		if err != nil {
			return false, err
		}

		if getRowsCount(rows) == 0 {
			inserted, err := insertDataValuesToDB(data, pManager.dataBase)
			if !inserted {
				return false, err
			}
		}
	}
	return true, nil
}

// Function for put model.Data entity to the db
func insertDataValuesToDB(data *model.Data, db *sql.DB) (bool, error) {
	_, err := db.Exec("INSERT INTO Pupils (id, pupil, establishment, class_num, letter) VALUES ($1, $2, $3, $4, $5)", data.ID, data.Pupil, data.Establishment, data.Class, data.Letter)
	return err == nil, err
}

// Function for determinate query rows size
func getRowsCount(rows *sql.Rows) uint64 {
	var size uint64 = 0
	for rows.Next() {
		size++
	}
	return size
}

// Function for get connection statement
func getOpenConnectionStatement(pInfo *ConnectionDTO) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s port=%s dbname=%s sslmode=%s",
		pInfo.HostName, pInfo.User, pInfo.Password, pInfo.Port, pInfo.DBName, "disable")
}

// Aux function for determ lower bound index in the sorted sequence
func lowerBound(sl []interface{}, left, right int, t time.Time) int {
	count := right - left
	var it int

	for count > 0 {
		it = left
		step := count / 2
		it += step

		curT := sl[it].(time.Time)
		comparisonRes := curT.Compare(t)
		if comparisonRes == -1 || comparisonRes == 0 {
			left++
			count -= step + 1
		} else {
			count = step
		}
	}

	return left
}

func TimesComparator(a, b interface{}) int {
	timeA := a.(time.Time)
	timeB := b.(time.Time)
	return timeA.Compare(timeB)
}
