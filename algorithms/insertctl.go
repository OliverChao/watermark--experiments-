package algorithms

import (
	"crypto/sha1"
	"hash"
	"math/big"
	"strconv"
)

import (
	"github.com/sirupsen/logrus"
)

type TestData struct {
	id     uint64
	name   string
	age    int
	weight float32
	height float32
}

//func Float32ToByte(float float32) []byte {
//	bits := math.Float32bits(float)
//	bytes := make([]byte, 4)
//	binary.LittleEndian.PutUint32(bytes, bits)
//
//	return bytes
//}

func InsertWaterMarking(ds []TestData, key []byte, gamma, nu, xi uint) {
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

	// todo: if there you uses database, you need to use tx.
	for r := range ds {
		firstHash.Reset()
		macHash.Reset()

		// H(K, r.P)
		p = []byte(strconv.FormatUint(ds[r].id, 10))
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

			//todo :  buf_chan  <- some_meta_data
			mark(bigFirstHash, attrIndex, bitIndex)
		}

		//	reset hash handle

	}
}

func getHashSum(h hash.Hash, bs ...[]byte) []byte {
	for _, b := range bs {
		h.Write(b)
	}
	return h.Sum(nil)
}

func mark(z big.Int, attrIndex int, bitIndex int) {
	flag := baseSelect(z, 2)
	if flag == 0 {
		logrus.Info("is even ...")
	} else {
		logrus.Info("is odd ...")
	}
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

//func Sha1Hash(bs []byte) []byte {
//	firstHash := sha1.New()
//	firstHash.Write(bs)
//	return firstHash.Sum(nil)
//}
