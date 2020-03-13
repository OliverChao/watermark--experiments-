package algorithms

import (
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGenerateToChan(t *testing.T) {

	model.FlagConfInit()
	service.ConnectDB()
	data := service.Student.GetBatchData(0, 1000)
	ch := make(chan *util.ChanMetaData, 50)
	var wg = sync.WaitGroup{}
	wg.Add(1)
	fmt.Println("go on.... core function")
	go GenerateToChan(data, ch, &wg)
	go service.Student.UpdateDB(ch)
	//go func() {
	//	con := 0
	//	for metaData := range ch {
	//		con++
	//		fmt.Println(metaData.ID,metaData.NewData)
	//	}
	//	fmt.Println("total ",con)
	//}()
	wg.Wait()
	close(ch)

	service.DisconnectDB()
	time.Sleep(time.Second)
}
