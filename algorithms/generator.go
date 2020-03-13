package algorithms

import (
	"WaterMasking/model"
	"WaterMasking/util"
	"crypto/sha1"
	"fmt"
	"hash"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"sync"
	"unsafe"
)

//var semaCon = 30
//var sema = make(chan struct{}, 30)

func GenerateToChan(ds []*model.Student, ch chan<- *util.ChanMetaData, wg *sync.WaitGroup) {

	defer wg.Done()

	key := []byte(model.Conf.Key)
	gamma := model.Conf.Gamma
	nu := model.Conf.Nu
	xi := model.Conf.Xi
	var (
		//the first hash handler
		firstHash = sha1.New()
		//the second hash handler
		macHash = sha1.New()

		//big hash number
		bigFirstHash = big.Int{}
		bigFMAC      = big.Int{}

		// p is the private key converted to []byte
		p []byte
		//fMAC []byte

		// attr_index and bit_index
		attrIndex, bitIndex int
	)

	for r := range ds {
		//logrus.Info("for range")
		firstHash.Reset()
		macHash.Reset()

		// H(K, r.P)
		p = []byte(strconv.FormatUint(uint64(ds[r].ID), 10))
		//firstHash.Write(key)
		//firstHash.Write(p)
		//firstSum := firstHash.Sum(nil)
		firstSum := getHashSum(firstHash, key, p)
		bigFirstHash.SetBytes(firstSum)

		//F(r.P)=H(K, H(K, r.P))
		//macHash.Write(key)
		//macHash.Write(firstSum)
		//fMAC := macHash.Sum(nil)
		fMAC := getHashSum(macHash, key, firstSum)
		bigFMAC.SetBytes(fMAC)
		//fMACHexString = hex.EncodeToString(fMAC)

		// core algorithm
		if TupleTest(bigFMAC, gamma) {
			attrIndex = AttrSelect(bigFMAC, nu)
			bitIndex = BitSelect(bigFMAC, xi)
			if flag, metaData := isNeedUpdate(bigFirstHash, ds[r], attrIndex, bitIndex); flag {
				ch <- metaData
			}
		}
	}
}

// if it returns false, it means there is no need to update.
func isNeedUpdate(bigFirstHash big.Int, d *model.Student, attrIndex int, bitIndex int) (bool, *util.ChanMetaData) {

	// reflect may impair the whole performance
	valueOf := reflect.ValueOf(d)
	fieldName := model.Conf.FiledNames[attrIndex]
	old := valueOf.Elem().FieldByName(fieldName).Float()

	var (
		flag    bool
		newData float64
	)
	f := baseSelect(bigFirstHash, 2)
	if f == 0 {
		flag, newData = cmpOldAndNew(old, bitIndex, '0')
	} else {
		flag, newData = cmpOldAndNew(old, bitIndex, '1')
	}

	if !flag {
		metaData := d.MetaData(fieldName, newData)
		return true, metaData
	}
	return false, nil

}

//If this function returns true, data[bitIndex] will be equal to the guess bit
//which means it is not necessary to update this piece of data.
//
func cmpOldAndNew(old float64, bitIndex int, guess byte) (bool, float64) {

	// change a data typed float to uint attribute to the IEEE 754
	ubits := math.Float64bits(old)
	// get the binary number typed string
	sbits := fmt.Sprintf("%b", ubits)
	bs := *(*[]byte)(unsafe.Pointer(&sbits))

	//i := 2
	//bs[len(bs)-1-i] = '0'
	//do not need to update this bit.
	if bs[len(bs)-1-bitIndex] == guess {
		return true, old
	}
	bs[len(bs)-1-bitIndex] = guess
	parseUint, _ := strconv.ParseUint(sbits, 2, 0)

	frombits := math.Float64frombits(parseUint)
	return false, frombits
}

func getHashSum(h hash.Hash, bs ...[]byte) []byte {
	for _, b := range bs {
		h.Write(b)
	}
	return h.Sum(nil)
}

func BitSelect(fmac big.Int, xi uint) int {
	return baseSelect(fmac, xi)
}

func TupleTest(fmac big.Int, gamma uint) bool {
	//logrus.Info(fmac)
	return baseSelect(fmac, gamma) == 0
}

func AttrSelect(fmac big.Int, nu uint) int {
	return baseSelect(fmac, nu)
}

func baseSelect(z big.Int, i uint) int {
	//logrus.Info("z is ", z.String(), " mod ", i)
	n := big.NewInt(int64(i))
	n.Mod(&z, n)
	u := n.Uint64()
	//logrus.Info("answer is ", u)
	return int(u)
}
