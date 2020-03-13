package algorithms

import (
	"WaterMasking/model"
	"WaterMasking/util"
	"crypto/sha1"
	"math/big"
	"strconv"
	"sync"
)

func VerifyData(ds []*model.Student, wg *sync.WaitGroup, ch chan<- util.VerifyResultData) {
	defer wg.Done()
	key := []byte(model.Conf.Key)
	gamma := model.Conf.Gamma
	nu := model.Conf.Nu
	xi := model.Conf.Xi
	var (
		totalCount, matchCount int
		//the first hash handler
		firstHash = sha1.New()
		//the second hash handler
		macHash = sha1.New()

		//big hash number
		bigFirstHash = big.Int{}
		bigFMAC      = big.Int{}

		// p is the private key converted to []byte
		p []byte

		// attr_index and bit_index
		attrIndex, bitIndex int
	)
	for r := range ds {
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
			totalCount++
			if flag, _ := isNeedUpdate(bigFirstHash, ds[r], attrIndex, bitIndex); !flag {
				matchCount = matchCount + 1
			}

		}
	}
	metaData := util.VerifyResultData{
		TotalCount: totalCount,
		MatchCount: matchCount,
	}
	ch <- metaData
}
