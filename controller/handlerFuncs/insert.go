package handlerFuncs

import (
	"WaterMasking/algorithms"
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func InsertWaterMarking(ctx *gin.Context) {
	ch := make(chan *util.ChanMetaData, 50)
	done := make(chan struct{})
	var wg = sync.WaitGroup{}
	wg.Add(2)

	var data []*model.Student

	switch service.SourceCache.IsCached() {
	case true:
		data = service.SourceCache.GetSourceData()
	case false:
		data = service.Student.GetAllData()
	}

	start := time.Now()
	batchSize := len(data) / 2
	go algorithms.GenerateToChan(data[:batchSize], ch, &wg)
	go algorithms.GenerateToChan(data[batchSize:], ch, &wg)
	go service.Student.UpdateDB(ch, done)
	go func() {
		wg.Wait()
		close(ch)
	}()
	<-done
	span := time.Since(start)

	logrus.Info("total time ", span)
	spanS := fmt.Sprintf("%v", span)
	ctx.JSON(200, map[string]interface{}{
		"total_time": spanS,
	})

	switch service.MarkedCache.IsCached() {
	case true:
	case false:
		service.MarkedCache.EnCached(service.Student.GetAllData())
	}

}
