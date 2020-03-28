package handlerFuncs

import (
	"WaterMasking/algorithms"
	"WaterMasking/model"
	"WaterMasking/service"
	"WaterMasking/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"math"
	"strconv"
	"sync"
	"unsafe"
)

func genDeleteData(d []*model.Student, p int) (ret []*model.Student) {
	num := len(d)
	ret = d[num/p:]
	return
}

func DeleteAttack(ctx *gin.Context) {
	if !service.MarkedCache.IsCached() {
		ctx.JSON(200, map[string]interface{}{
			"code": -1,
			"msg":  "have not inserted watermark",
		})
		return
	}

	rets := make([]util.VerifyResultData, 4)
	rets = rets[:0]
	for i := 5; i >= 2; i-- {
		data := genDeleteData(service.MarkedCache.GetSourceData(), i)
		ret := _verify(data)
		rets = append(rets, ret)
	}

	ctx.JSON(200, map[string]interface{}{
		"data": rets,
	})

}

func AddMoreAttack(ctx *gin.Context) {

	if !service.MarkedCache.IsCached() {
		ctx.JSON(200, map[string]interface{}{
			"code": -1,
			"msg":  "have not inserted watermark",
		})
		return
	}

	alpha := ctx.DefaultQuery("alpha", "0.9")
	// error handler is eliminated.
	// please set correct alpha ((0.5,1))
	f, _ := strconv.ParseFloat(alpha, 64)
	M := len(service.SourceCache.GetSourceData())
	num := int(((1-f)*float64(M))/(f-0.5)) + 1

	newData := service.SourceCache.GetSourceData()[:num]

	attackData := append(service.MarkedCache.GetSourceData(), newData...)

	ret := _verify(attackData)

	currentAlpha := float64(ret.MatchCount) / float64(ret.TotalCount)
	logrus.Info("alpha = ", currentAlpha)

	ctx.JSON(200, map[string]interface{}{
		"before_add":    M,
		"after_add":     M + num,
		"add_num":       num,
		"data":          ret,
		"current_alpha": currentAlpha,
	})

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

func ReverseAttack(ctx *gin.Context) {
	if !service.MarkedCache.IsCached() {
		ctx.JSON(200, map[string]interface{}{
			"code": -1,
			"msg":  "have not inserted watermark",
		})
		return
	}
	MS := ctx.DefaultPostForm("number", "100000")
	alphaS := ctx.DefaultPostForm("alpha", "0.9")

	MI, _ := strconv.ParseInt(MS, 10, 0)
	alpha, _ := strconv.ParseFloat(alphaS, 64)
	//num :=
	baseData := service.MarkedCache.GetSourceData()
	M := int(math.Min(float64(MI), float64(len(baseData))))

	attackNum := int(float64(M)*(1-alpha)) + 1
	data := genBackData(baseData[:M], attackNum)

	ret := _verify(data)

	currentAlpha := float64(ret.MatchCount) / float64(ret.TotalCount)
	logrus.Info("alpha = ", currentAlpha)

	ctx.JSON(200, map[string]interface{}{
		"reverse":       attackNum,
		"data":          ret,
		"current_alpha": currentAlpha,
	})

}

func genBackData(data []*model.Student, num int) []*model.Student {
	d := make([]*model.Student, len(data))
	d = d[:0]
	for i := 0; i < num; i++ {
		m := *data[i]
		reverseBits(&m)
		d = append(d, &m)
	}
	d = append(d, data[num:]...)
	return d
}

func reverseBits(s *model.Student) {
	//s.Score1
	b1 := getBits(s.Score1)
	b2 := getBits(s.Score2)
	b3 := getBits(s.Score3)
	b4 := getBits(s.Score4)
	b5 := getBits(s.Score5)

	xi := int(model.Conf.Xi)
	b1 = reverse(b1, xi)

	b2 = reverse(b2, xi)
	b3 = reverse(b3, xi)
	b4 = reverse(b4, xi)
	b5 = reverse(b5, xi)
	s1 := *(*string)(unsafe.Pointer(&b1))
	s2 := *(*string)(unsafe.Pointer(&b2))
	s3 := *(*string)(unsafe.Pointer(&b3))
	s4 := *(*string)(unsafe.Pointer(&b4))
	s5 := *(*string)(unsafe.Pointer(&b5))

	parseUint1, _ := strconv.ParseUint(s1, 2, 0)
	parseUint2, _ := strconv.ParseUint(s2, 2, 0)
	parseUint3, _ := strconv.ParseUint(s3, 2, 0)
	parseUint4, _ := strconv.ParseUint(s4, 2, 0)
	parseUint5, _ := strconv.ParseUint(s5, 2, 0)

	frombits1 := math.Float64frombits(parseUint1)
	s.Score1 = frombits1
	frombits2 := math.Float64frombits(parseUint2)
	s.Score2 = frombits2
	frombits3 := math.Float64frombits(parseUint3)
	s.Score3 = frombits3
	frombits4 := math.Float64frombits(parseUint4)
	s.Score4 = frombits4
	frombits5 := math.Float64frombits(parseUint5)
	s.Score5 = frombits5

}

func getBits(f float64) []byte {
	ubits := math.Float64bits(f)
	// get the binary number typed string
	sbits := fmt.Sprintf("%b", ubits)
	bs := *(*[]byte)(unsafe.Pointer(&sbits))
	return bs
}

func reverse(bs []byte, xi int) []byte {
	_ = bs[len(bs)-1] // bounds check elimination

	for i := len(bs) - xi; i <= len(bs)-1; i++ {
		//bs[i] = bs[i+1]
		if bs[i] == '1' {
			bs[i] = '0'
		} else if bs[i] == '0' {
			bs[i] = '1'
		}
	}
	return bs
}
