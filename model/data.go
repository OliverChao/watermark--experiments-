package model

import "WaterMasking/util"

// data struct
type Student struct {
	//gorm.Model
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"type:varchar(20)"`
	School string `gorm:"size:100"`
	Score1 float64
	Score2 float64
	Score3 float64
	Score4 float64
	Score5 float64
}

type StudentBack struct {
	//gorm.Model
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"type:varchar(20)"`
	School string `gorm:"size:100"`
	Score1 float64
	Score2 float64
	Score3 float64
	Score4 float64
	Score5 float64
}

//
func (s *Student) MetaData(col string, newData float64) *util.ChanMetaData {
	ret := &util.ChanMetaData{
		ID:      s.ID,
		Col:     col,
		NewData: newData,
	}
	return ret
}
