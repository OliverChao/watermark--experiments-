package handlerFuncs

import (
	"WaterMasking/model"
	"WaterMasking/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"sync"
	"time"
	"unsafe"
)

func BackupData(ctx *gin.Context) {
	start := time.Now()
	var wrWg = sync.WaitGroup{}
	data := service.SourceCache.GetSourceData()
	service.StudentBack.ClearTable(&model.StudentBack{})
	dataBack := *(*[]*model.StudentBack)(unsafe.Pointer(&data))
	grNum := int(math.Min(float64(model.Conf.Gamma), 3))
	//backCh := make(chan )
	batchSize := len(data) / grNum

	wrWg.Add(grNum)
	for i := 0; i < grNum-1; i++ {
		go service.StudentBack.AsyncWriteStudentBacks(dataBack[batchSize*i:batchSize*(i+1)], &wrWg)
	}
	go service.StudentBack.AsyncWriteStudentBacks(dataBack[batchSize*(grNum-1):], &wrWg)

	wrWg.Wait()

	span := time.Since(start)
	spanS := fmt.Sprintf("%v", span)
	logrus.Info("total time ", span)
	ctx.JSON(200, map[string]interface{}{
		"total_time": spanS,
	})

}
