package service

import (
	"WaterMasking/model"
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
