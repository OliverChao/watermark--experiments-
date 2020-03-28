package attack

import (
	"WaterMasking/model"
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

func ReverseAttack(data []*model.Student) []*model.Student {
	for i := 0; i < 10000; i++ {
		reverseBits(data[i])
	}
	return data
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
