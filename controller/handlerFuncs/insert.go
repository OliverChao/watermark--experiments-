package handlerFuncs

import (
	"WaterMasking/algorithms"
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"math"
	"sync"
	"time"
)

type insertParameters struct {
	Table string `form:"table" json:"table" binding:"required"`
	Key   string `form:"key" json:"key" binding:"required"`
	Gamma uint   `form:"gamma" json:"gamma" binding:"required"`
	Nu    uint   `form:"nu" json:"nu" binding:"required"`
	Xi    uint   `form:"xi" json:"xi" binding:"required"`
}

func InsertWaterMarking(ctx *gin.Context) {

	//table := ctx.DefaultPostForm()
	var params insertParameters
	if err := ctx.MustBindWith(&params, binding.FormPost); err != nil {
		ctx.String(200, err.Error())
		return
	}
	//ctx.JSON(200, params)
	//return
	logrus.Info(params)
	updateConfByInserted(params)

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

	service.ResumeController.UnResumed()
	logrus.Info("total time ", span)
	spanS := fmt.Sprintf("%v", span)
	ctx.JSON(200, map[string]interface{}{
		"total_time": spanS,
	})

	service.MarkedCache.EnCached(service.Student.GetAllData())

}

func updateConfByInserted(i insertParameters) {
	model.Conf.Gamma = i.Gamma
	defaultFileds := []string{"Score1", "Score2", "Score3", "Score4", "Score5"}
	if i.Nu > 0 {
		model.Conf.Nu = uint(math.Min(float64(5), float64(i.Nu)))
	}
	model.Conf.FiledNames = defaultFileds[:model.Conf.Nu]
	model.Conf.Xi = i.Xi
	model.Conf.Key = i.Key

}
