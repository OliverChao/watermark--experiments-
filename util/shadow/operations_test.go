package shadow

import (
	"fmt"
	"testing"
)

func TestRsaDecrypt(t *testing.T) {
	//RegisterCryptoSystem()
	b, _ := RsaEncrypt([]byte("将进酒，杯莫停。"))
	fmt.Printf("%x\n", b)
	b2, _ := RsaDecrypt(b)
	fmt.Println(string(b2))
}

func TestRsaEncrypt(t *testing.T) {
	//RegisterCryptoSystem()
	b, _ := RsaEncrypt([]byte("将进酒，杯莫停。"))
	fmt.Printf("%x\n", b)
}

func TestRsaSign(t *testing.T) {
	//RegisterCryptoSystem()
	sign, e := RsaSign([]byte("oliver loves annabelle"))
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Printf("%x\n", sign)

	}

}

func TestRsaVerify(t *testing.T) {
	msg := []byte("oliver loves annabelle")
	sign, e := RsaSign(msg)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Printf("%x\n", sign)
	}
	e = RsaVerify(msg, sign)
	if e != nil {
		fmt.Println("verify successfully ~~")
	}
}
