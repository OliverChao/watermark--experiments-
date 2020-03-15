package main

import (
	"WaterMasking/algorithms"
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"sync"
	"time"
	"unsafe"
)

func init() {
	model.FlagConfInit()
}

func main() {
	//fmt.Println(model.Conf.FiledNames, model.Conf.Nu)
	prgStart := time.Now()
	service.ConnectDB()
	if model.Conf.ExecMode == "verify" {
		verify()
	} else {
		insert()
	}

	//service.DestroyAll()

	service.DisconnectDB()
	fmt.Println("total: ", time.Since(prgStart))
}

func verify() {
	start := time.Now()
	data := service.Student.GetAllData()
	ch := make(chan util.VerifyResultData, 1)
	var wg = sync.WaitGroup{}
	wg.Add(1)
	algorithms.VerifyData(data, &wg, ch)
	wg.Wait()
	fmt.Println(<-ch)
	fmt.Println("verify: ", time.Since(start))
}

func insert() {
	ch := make(chan *util.ChanMetaData, 30)
	done := make(chan struct{})
	var wg = sync.WaitGroup{}

	wg.Add(2)
	start := time.Now()
	data := service.Student.GetAllData()
	batchSize := len(data) / 2
	go algorithms.GenerateToChan(data[:batchSize], ch, &wg)
	go algorithms.GenerateToChan(data[batchSize:], ch, &wg)
	go service.Student.UpdateDB(ch, done)
	if model.Conf.BakeUp {
		//service.DropTableIfExists(&model.StudentBack{})
		//db.AutoMigrate(&model.StudentBack{})
		service.StudentBack.ClearTable(&model.StudentBack{})
		dataBack := *(*[]*model.StudentBack)(unsafe.Pointer(&data))
		grNum := int(math.Min(float64(model.Conf.Gamma), 3))
		//backCh := make(chan )
		batchSize := len(data) / grNum
		wg.Add(grNum)
		for i := 0; i < grNum-1; i++ {
			go service.StudentBack.AsyncWriteStudentBacks(dataBack[batchSize*i:batchSize*(i+1)], &wg)
		}
		go service.StudentBack.AsyncWriteStudentBacks(dataBack[batchSize*grNum:], &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	<-done
	logrus.Info("total:", time.Since(start))
	//logrus.Info("read data: ", time.Since(readStart), "s")
}
