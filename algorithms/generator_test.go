package algorithms

import (
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"sync"
	"testing"
)

func TestGenerateToChan(t *testing.T) {

	model.FlagConfInit()
	service.ConnectDB()
	data := service.Student.GetBatchData(0, 100000)
	ch := make(chan *util.ChanMetaData, 30)
	done := make(chan struct{})
	var wg = sync.WaitGroup{}
	wg.Add(2)
	fmt.Println("go on.... core function")
	//go GenerateToChan(data[:1500], ch, &wg)
	go GenerateToChan(data[:50000], ch, &wg)
	go GenerateToChan(data[50000:], ch, &wg)

	go service.Student.UpdateDB(ch, done)
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
	//time.Sleep(time.Second)
	<-done
}

func Test_cmpOldAndNew(t *testing.T) {
	old := 3.14
	//u := *(*[]byte)(unsafe.Pointer(&old))
	//fmt.Println(u)
	ubits := math.Float64bits(old)
	logrus.Info(ubits)
	// get the binary number typed string
	sbits := fmt.Sprintf("%b", ubits)
	logrus.Println(sbits)
	//bs := *(*[]byte)(unsafe.Pointer(&sbits))
	//logrus.Info(string(bs))
}
