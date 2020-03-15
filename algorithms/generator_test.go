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
	data := service.Student.GetBatchData(0, 100000)
	ch := make(chan *util.ChanMetaData, 30)
	var wg = sync.WaitGroup{}
	wg.Add(2)
	fmt.Println("go on.... core function")
	//go GenerateToChan(data[:1500], ch, &wg)
	go GenerateToChan(data[:50000], ch, &wg)
	go GenerateToChan(data[50000:], ch, &wg)
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
