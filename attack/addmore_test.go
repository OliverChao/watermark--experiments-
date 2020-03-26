package attack

import (
	"WaterMasking/algorithms"
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"sync"
	"testing"
)

func TestAddMoreAttack(t *testing.T) {
	model.FlagConfInit()
	service.ConnectDB()
	//dataBack := service.StudentBack.GetAllData()

	dataBack := service.GetAllDataByTableName("student_backs")
	dataSource := service.Student.GetBatchData(0, 25000)

	for i := range dataSource {
		dataSource[i].ID = uint(100000 + i + 1)
	}

	fmt.Println(len(dataSource), len(dataBack))
	dataBase := make([]*model.Student, len(dataBack)+len(dataSource)+1)
	fmt.Println(len(dataBase))
	data := dataBase[:0]
	data = append(data, dataBack...)
	data = append(data, dataSource...)

	fmt.Println(len(data), data[len(data)-1].ID)
	ch := make(chan util.VerifyResultData, 1)
	var wg = sync.WaitGroup{}
	wg.Add(1)
	//algorithms.VerifyData(data[:len(dataBack)], &wg, ch)
	//algorithms.VerifyData(data[len(dataBack):], &wg, ch)
	algorithms.VerifyData(data, &wg, ch)
	wg.Wait()
	resultData := <-ch
	//fmt.Println(<-ch)
	fmt.Println(resultData)
	fmt.Println(float64(resultData.MatchCount) / float64(resultData.TotalCount))
	service.DisconnectDB()

}
