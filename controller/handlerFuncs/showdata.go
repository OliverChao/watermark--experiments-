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
	offsetS := c.DefaultQuery("offset", "1")
	countS := c.DefaultQuery("count", "500")
	back := c.DefaultQuery("data", "source")
	offsetI, _ := strconv.Atoi(offsetS)
	countI, _ := strconv.Atoi(countS)
	if offsetI < 1 {
		logrus.Error("invalid offset, set offset=1")
		offsetI = 1
	}
	if countI < 1 {
		logrus.Error("invalid count, set count=100")
		countI = 100
	}

	var data []*model.Student
	switch back {
	case "source":
		data = service.SourceCache.GetSourceData()[offsetI-1 : offsetI+countI-1]
	case "back":
		if !service.MarkedCache.IsCached() {
			c.String(200, "have not inserted watermark")
			return
		}
		data = service.MarkedCache.GetSourceData()[offsetI-1 : offsetI+countI-1]
	}
	//c.JSON(200, map[string]interface{}{
	//	"data": data,
	//})
	c.HTML(http.StatusOK, "showdata.html", gin.H{
		"data":  data,
		"table": []string{"students"},
	})
}

func ShowDataPost(c *gin.Context) {
	offsetS := c.DefaultPostForm("offset", "0")
	countS := c.DefaultPostForm("count", "500")
	back := c.DefaultPostForm("data", "source")
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
	switch back {
	case "source":
		data = service.SourceCache.GetSourceData()[offsetI : offsetI+countI]
	case "back":
		if !service.MarkedCache.IsCached() {
			c.String(200, "have not inserted watermark")
			return
		}
		data = service.MarkedCache.GetSourceData()[offsetI-1 : offsetI+countI]
	}
	//c.JSON(200, map[string]interface{}{
	//	"data": data,
	//})
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
