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

func DestroyAll() {
	db.DropTableIfExists(&model.Student{})
	db.AutoMigrate(&model.Student{})
}

func ChangeTableName(s, n string) {
	sql := (" " + "alter table " + s + " rename to " + n + ";")[1:]
	db.Exec(sql)
}
