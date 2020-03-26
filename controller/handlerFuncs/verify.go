package handlerFuncs

import (
	"WaterMasking/algorithms"
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"github.com/b3log/gulu"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func VerifyWaterMarking(ctx *gin.Context) {
	result := gulu.Ret.NewResult()
	defer ctx.JSON(200, result)

	var data []*model.Student
	query := ctx.DefaultQuery("data", "source")

	switch query {
	case "source":
		data = service.SourceCache.GetSourceData()
	case "back":
		if !service.MarkedCache.IsCached() {
			result.Msg = "have not inserted watermark"
			result.Code = -1
			return
		}
		data = service.MarkedCache.GetSourceData()
	}

	start := time.Now()
	ch := make(chan util.VerifyResultData, 1)
	var wg = sync.WaitGroup{}
	wg.Add(1)
	algorithms.VerifyData(data, &wg, ch)
	wg.Wait()
	//fmt.Println(<-ch)
	ret := <-ch
	result.Data = ret
	result.Msg = fmt.Sprintf("%v", time.Since(start))

	logrus.Info("a = ", float64(ret.MatchCount)/float64(ret.TotalCount))
}
