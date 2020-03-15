package service

import (
	"WaterMasking/model"
	"testing"
)

func TestChangeTableName(t *testing.T) {
	model.FlagConfInit()
	ConnectDB()

	ChangeTableName("students", "students_back")
	DisconnectDB()
}
