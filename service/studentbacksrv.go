package service

import (
	"WaterMasking/model"
	"WaterMasking/util"
	"github.com/sirupsen/logrus"
	"sync"
)

// student data operating service
type studentBackSrv struct {
	mutex *sync.Mutex
}

var StudentBack = &studentBackSrv{mutex: &sync.Mutex{}}

func (srv *studentBackSrv) GetBatchData(offset, limit int) (ret []*model.StudentBack) {
	if limit < 0 {
		return nil
	}
	if offset < 0 {
		offset = 0
	}

	if err := db.Model(&model.StudentBack{}).Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		logrus.Error("get batch data error", err)
	}
	return
}

func (srv *studentBackSrv) GetAllData() (ret []*model.StudentBack) {
	if err := db.Model(&model.StudentBack{}).Find(&ret).Error; err != nil {
		logrus.Error(err)
	}
	return
}

//func (srv *studentBackSrv) UpdateDB
// use this method like the following
// go service.StudentBack.UpdateDB(ch)
func (srv *studentBackSrv) UpdateDB(ch <-chan *util.ChanMetaData) (err error) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
			logrus.Info("tx commit successfully")
		} else {
			tx.Rollback()
			logrus.Error(err, "rollback...")
		}
	}()

Loop:
	for {
		select {
		case metadata, ok := <-ch:
			if !ok {
				break Loop
			}
			tx.Model(&model.StudentBack{ID: metadata.ID}).Update(metadata.Col, metadata.NewData)
		}
	}

	return nil
}

func (srv *studentBackSrv) SyncWriteOneStudentBack(s *model.StudentBack) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()
	db.Create(&s)
}

// async code
// the `id` must be controlled correctly by user.
// or will arise error, say, `Duplicate entry [num] for key 'PRIMARY'`
func (srv *studentBackSrv) AsyncWriteStudentBacks(ss []*model.StudentBack, wg *sync.WaitGroup) {
	writeStudentsBack(ss, wg)
}

//// sync transaction lock
//func (srv *studentBackSrv) SyncWriteStudentBacks(ss []*model.StudentBack) {
//	srv.mutex.Lock()
//	defer srv.mutex.Unlock()
//	writeStudentBacks(ss)
//}

// underlying functions for writing model.StudentBack to table
func writeStudentsBack(ss []*model.StudentBack, wg *sync.WaitGroup) {
	defer wg.Done()

	tx := db.Begin()
	var err error
	defer func() {
		if nil != err {
			logrus.Error(err)
			tx.Rollback()
		} else {
			logrus.Info("commit successfully")
			tx.Commit()
		}
	}()
	for i := range ss {
		tx.Create(&ss[i])
	}
}

func (srv *studentBackSrv) ClearTable(values ...interface{}) {
	db.DropTableIfExists(values...)
	db.AutoMigrate(values...)
}
