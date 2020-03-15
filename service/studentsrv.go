package service

import (
	"WaterMasking/model"
	"WaterMasking/util"
	"github.com/sirupsen/logrus"
	"sync"
)

// student data operating service
type studentSrv struct {
	mutex *sync.Mutex
}

var Student = &studentSrv{mutex: &sync.Mutex{}}

func (srv *studentSrv) GetBatchData(offset, limit int) (ret []*model.Student) {
	if limit < 0 {
		return nil
	}
	if offset < 0 {
		offset = 0
	}

	if err := db.Model(&model.Student{}).Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		logrus.Error("get batch data error", err)
	}
	return
}

func (srv *studentSrv) GetAllData() (ret []*model.Student) {
	if err := db.Model(&model.Student{}).Find(&ret).Error; err != nil {
		logrus.Error(err)
	}
	return
}

//func (srv *studentSrv) UpdateDB
// use this method like the following
// go service.Student.UpdateDB(ch,done)
// use  `<- done` to check weather or not the function is completed
func (srv *studentSrv) UpdateDB(ch <-chan *util.ChanMetaData, done chan<- struct{}) (err error) {
	defer func() {
		done <- struct{}{}
	}()

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
			tx.Model(&model.Student{ID: metadata.ID}).Update(metadata.Col, metadata.NewData)
		}
	}

	return nil
}

func (srv *studentSrv) SyncWriteOneStudent(s *model.Student) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()
	db.Create(&s)
}

// async code
// the `id` must be controlled correctly by user.
// or will arise error, say, `Duplicate entry [num] for key 'PRIMARY'`
func (srv *studentSrv) AsyncWriteStudents(ss []*model.Student) {
	writeStudents(ss)
}

// sync transaction lock
func (srv *studentSrv) SyncWriteStudents(ss []*model.Student) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()
	writeStudents(ss)
}

// underlying functions for writing model.Student to table
func writeStudents(ss []*model.Student) {
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
