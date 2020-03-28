package attack

import (
	"WaterMasking/algorithms"
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
)

func Test_reverse(t *testing.T) {
	bs := []byte("oliver loves oliver")
	reverse(bs, 6)
	fmt.Println(string(bs))
}

func TestReverseAttack(t *testing.T) {
	model.FlagConfInit()
	service.ConnectDB()
	data := service.Student.GetAllData()[:]
	data = ReverseAttack(data)
	ret := _verify(data)
	logrus.Info(ret)
	logrus.Info(float64(ret.MatchCount) / float64(ret.TotalCount))

	service.DisconnectDB()
}

func _verify(data []*model.Student) util.VerifyResultData {
	ch := make(chan util.VerifyResultData, 1)
	var wg = sync.WaitGroup{}
	wg.Add(1)
	algorithms.VerifyData(data, &wg, ch)
	wg.Wait()
	//fmt.Println(<-ch)
	ret := <-ch
	return ret
}
