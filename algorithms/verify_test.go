package algorithms

import (
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"sync"
	"testing"
)

func TestVerifyData(t *testing.T) {
	model.FlagConfInit()
	service.ConnectDB()
	data := service.Student.GetBatchData(0, 1000)
	ch := make(chan util.VerifyResultData, 1)
	var wg = sync.WaitGroup{}
	wg.Add(1)
	fmt.Println("go on.... core function")
	VerifyData(data, &wg, ch)
	wg.Wait()
	fmt.Println(<-ch)
	service.DisconnectDB()
}
