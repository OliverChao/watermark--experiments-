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

//func (srv *studentSrv) UpdateDB
// use this method like the following
// go service.Student.UpdateDB(ch)
func (srv *studentSrv) UpdateDB(ch <-chan *util.ChanMetaData) (err error) {
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
