package service

import (
	"WaterMasking/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"time"
)

var db *gorm.DB

func ConnectDB() {
	var err error
	db, err = gorm.Open("mysql", model.Conf.Server)

	logrus.Info("[mysql dial]: ", model.Conf.Server)
	if err != nil {
		logrus.Fatal("Cannot connect to the database", err)
		return
	}

	logrus.Info("Connect to mysql successfully")

	db.AutoMigrate(&model.Student{})
	// some database connection configurations
	db.DB().SetMaxIdleConns(15)
	db.DB().SetMaxOpenConns(55)
	db.DB().SetConnMaxLifetime(5 * time.Minute)
}

func DisconnectDB() {
	if err := db.Close(); nil != err {
		logrus.Error("Disconnect from database failed: " + err.Error())
	}
	logrus.Info("Disconnect from database successfully")
}

func InitSourceCacheData() {
	SourceCache.EnCached(Student.GetAllData())
}

func DestroyAll() {
	db.DropTableIfExists(&model.Student{})
	db.AutoMigrate(&model.Student{})
}

func DropTableIfExists(values ...interface{}) {
	db.DropTableIfExists(values...)
}

func ChangeTableName(s, n string) {
	sql := (" " + "alter table " + s + " rename to " + n + ";")[1:]
	db.Exec(sql)
}

func ExchangeTableName() {
	sql1 := "alter table students rename to students_tmp"
	sql2 := "alter table student_backs rename to students"
	sql3 := "alter table students_tmp rename to student_backs"
	db.Exec(sql1)
	db.Exec(sql2)
	db.Exec(sql3)
}

func GetBatchDataByTableName(tname string, offset, limit int) (ret []*model.Student) {
	if limit < 0 {
		return nil
	}
	if offset < 0 {
		offset = 0
	}
	switch tname {
	case "students":
	case "student_backs":
	default:
		logrus.Error("do not support this table name")
		return nil
	}
	if err := db.Table(tname).Offset(offset).Limit(limit).Find(&ret).Error; err != nil {
		logrus.Error("get batch data error", err)
	}
	return
}

func GetAllDataByTableName(tname string) (ret []*model.Student) {
	if err := db.Table(tname).Find(&ret).Error; err != nil {
		logrus.Error(err)
	}
	return
}
