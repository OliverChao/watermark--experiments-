package service

import (
	"WaterMasking/model"
	"fmt"
	"sync"
	"testing"
	"unsafe"
)

// adding 100000 consumes time
func TestChangeTableName(t *testing.T) {
	model.FlagConfInit()
	ConnectDB()

	//ChangeTableName("students", "students_back")
	sData := Student.GetBatchData(0, 8000)
	//sData := Student.GetAllData()
	data := *(*[]*model.StudentBack)(unsafe.Pointer(&sData))
	fmt.Println(data[0].ID, data[0].Name, data[0].Score5)

	//db.DropTableIfExists(&model.StudentBack{})
	//db.AutoMigrate(&model.StudentBack{})
	StudentBack.ClearTable(&model.StudentBack{})
	var wg sync.WaitGroup
	n := 1
	wg.Add(n)
	batchSize := len(data) / n
	go writeStudentsBack(data[:batchSize], &wg)
	//go writeStudentsBack(data[batchSize:batchSize*2], &wg)
	//go writeStudentsBack(data[batchSize*2:batchSize*3], &wg)
	//go writeStudentsBack(data[batchSize*3:], &wg)
	//go writeStudentsBack(data[batchSize*4:], &wg)
	wg.Wait()
	DisconnectDB()
}

// updating 100000 consumes time
func TestChangeTableName2(t *testing.T) {
	model.FlagConfInit()
	ConnectDB()

	//ChangeTableName("students", "students_back")
	//sData := Student.GetBatchData(0, 1000)
	sData := Student.GetAllData()
	data := *(*[]*model.StudentBack)(unsafe.Pointer(&sData))
	fmt.Println(data[0].ID, data[0].Name, data[0].Score5)

	tx := db.Begin()
	for r := range data {
		tx.Model(&data[r]).Update("Score1", 1000.0)
	}
	tx.Commit()
	//db.DropTableIfExists(&model.StudentBack{})

	//db.AutoMigrate(&model.StudentBack{})
	//writeStudentsBack(data)
	//if err := db.Model(&model.StudentBack{}).Create(data[0]); err != nil {
	//	logrus.Error(err)
	//}
	//DestroyAll()
	DisconnectDB()
}

func TestDestroyAll(t *testing.T) {
	model.FlagConfInit()
	ConnectDB()
	DestroyAll()
	DisconnectDB()
}

//func writeStudentsBack(ss []*model.StudentBack, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	tx := db.Begin()
//	var err error
//	defer func() {
//		if nil != err {
//			logrus.Error(err)
//			tx.Rollback()
//		} else {
//			logrus.Info("commit successfully")
//			tx.Commit()
//		}
//	}()
//	for i := range ss {
//		tx.Create(&ss[i])
//	}
//}
