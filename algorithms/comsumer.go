package algorithms

import (
	"WaterMasking/model"
	"WaterMasking/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// This is logic code in order to understand.
//The actual execution code is written in studentsrv.go.
func Consumer(db *gorm.DB, ch <-chan *util.ChanMetaData) (err error) {
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
