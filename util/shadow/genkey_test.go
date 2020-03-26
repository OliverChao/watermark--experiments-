package shadow

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestGenRsaKey(t *testing.T) {
	var length int = 1024
	err := GenRsaKey(length)
	if err != nil {
		logrus.Error("gen rsa err")
	} else {
		logrus.Info("ok")
	}

}

func TestSomeFunction(t *testing.T) {
	//fmt.Println(os.Getwd())
	RegisterCryptoSystem()
	//ForceRegisterCryptoSystem()
	fmt.Println(cap(privateKey))
	fmt.Println(cap(publicKey))
	fmt.Println(len(privateKey))
	fmt.Println(len(publicKey))
}
