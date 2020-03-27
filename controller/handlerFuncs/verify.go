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

type VerifyShowData struct {
	util.VerifyResultData
	SignificantLevel float64
	Msg              string
}

func VerifyWaterMarking(ctx *gin.Context) {
	var tables []string
	switch service.MarkedCache.IsCached() {
	case true:
		tables = []string{"students", "student_backs"}
	case false:
		tables = []string{"students"}
	}
	result := VerifyShowData{}
	//defer ctx.HTML(200, "verify.html", gin.H{
	//	"table":  tables,
	//	"result": result,
	//})

	var data []*model.Student
	query := ctx.DefaultPostForm("table", "students")

	switch query {
	case "students":
		data = service.SourceCache.GetSourceData()
	case "student_backs":
		if !service.MarkedCache.IsCached() {
			result.Msg = "[error] have not inserted watermark"
			ctx.HTML(200, "verify.html", gin.H{
				"table":     tables,
				"result":    result,
				"testTable": query,
			})
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
	result.VerifyResultData = ret
	result.Msg = fmt.Sprintf("[successfully verify] spend %v", time.Since(start))
	result.SignificantLevel = float64(result.MatchCount) / float64(result.TotalCount)
	logrus.Info("significant lever = ", result.SignificantLevel)

	ctx.HTML(200, "verify.html", gin.H{
		"table":     tables,
		"result":    result,
		"testTable": query,
	})
}
