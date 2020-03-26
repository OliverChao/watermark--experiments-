package shadow

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

const (
	priKeyFile = "rsa_private_key.pem"
	pubKeyFile = "rsa_public_key.pem"
)

var (
	pathPri    string
	pathPub    string
	privateKey []byte
	publicKey  []byte
)

func init() {
	//fmt.Println("this is init function")
	dir, _ := os.Getwd()
	pathPri = dir + "/" + priKeyFile
	pathPub = dir + "/" + pubKeyFile

	RegisterCryptoSystem()
}

func RegisterCryptoSystem() {
	genCryptoFile()
	privateKey, _ = ioutil.ReadFile(pathPri)
	publicKey, _ = ioutil.ReadFile(pathPub)
}

func ForceRegisterCryptoSystem() {
	privateKey = make([]byte, 0)
	publicKey = make([]byte, 0)
	forceGenCryptoFiles()
	privateKey, _ = ioutil.ReadFile(pathPri)
	publicKey, _ = ioutil.ReadFile(pathPub)
}

func genCryptoFile() {

	if _, e := os.Stat(pathPri); e != nil {
		logrus.Info("there is no pri key")
		forceGenCryptoFiles()
		return
	}
	if _, e := os.Stat(pathPub); e != nil {
		logrus.Info("there is no pub key")
		forceGenCryptoFiles()
		return
	}
	logrus.Info("having rsa files successfully")

}

func forceGenCryptoFiles() {
	var rsaLength int = 1024
	e := GenRsaKey(rsaLength)
	if e != nil {
		logrus.Error("[ForceCryptoSystem] gen rsa error")
	}
}

func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	privFile, err := os.Create(priKeyFile)
	if err != nil {
		return err
	}
	defer privFile.Close()

	err = pem.Encode(privFile, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubFile, err := os.Create(pubKeyFile)
	if err != nil {
		return err
	}
	defer pubFile.Close()
	err = pem.Encode(pubFile, block)
	if err != nil {
		return err
	}
	return nil
}
