package handlerFuncs

import (
	"WaterMasking/model"
	"WaterMasking/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func ShowData(c *gin.Context) {
	offsetS := c.DefaultQuery("offset", "0")
	countS := c.DefaultQuery("count", "500")
	offsetI, _ := strconv.Atoi(offsetS)
	countI, _ := strconv.Atoi(countS)
	if offsetI < 0 {
		logrus.Error("invalid offset, set offset=1")
		offsetI = 0
	}
	if countI < 1 {
		logrus.Error("invalid count, set count=100")
		countI = 500
	}

	var data []*model.Student
	data = service.SourceCache.GetSourceData()[offsetI : offsetI+countI]
	var tables []string

	switch service.MarkedCache.IsCached() {
	case true:
		tables = []string{"students", "student_backs"}
	case false:
		tables = []string{"students"}
	}

	c.HTML(http.StatusOK, "showdata.html", gin.H{
		"data":  data,
		"table": tables,
	})
}

func ShowDataPost(c *gin.Context) {
	offsetS := c.DefaultPostForm("offset", "0")
	countS := c.DefaultPostForm("count", "500")
	table := c.DefaultPostForm("table", "students")
	offsetI, _ := strconv.Atoi(offsetS)
	countI, _ := strconv.Atoi(countS)

	total := len(service.SourceCache.GetSourceData())
	if (offsetI < 0 || offsetI > total) || (countI < 1 || countI > total) || (offsetI+countI > total) {
		logrus.Error("invalid offset, set offset=1")
		//offsetI = 1
		c.HTML(200, "showdataHint.html", gin.H{
			"msg": fmt.Sprintf("invalid arguments!! total_data : %d", total),
		})
		c.HTML(http.StatusOK, "showdata.html", gin.H{
			//"data": data,
			"table": []string{"students"},
		})
		return
	}

	var data []*model.Student

	switch table {
	case "students":
		data = service.SourceCache.GetSourceData()[offsetI : offsetI+countI]
	case "student_backs":
		if !service.MarkedCache.IsCached() {
			c.String(200, "have not inserted watermark")
			return
		}
		data = service.MarkedCache.GetSourceData()[offsetI : offsetI+countI]
	}

	var tables []string
	switch service.MarkedCache.IsCached() {
	case true:
		tables = []string{"students", "student_backs"}
	case false:
		tables = []string{"students"}
	}
	c.HTML(http.StatusOK, "showdata.html", gin.H{
		"data":  data,
		"table": tables,
	})

}
