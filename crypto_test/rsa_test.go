package cryptotest_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io/ioutil"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestRsa(t *testing.T) {

	var (
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		msg        []byte = []byte("hello world")
	)

	publicKeyByte, err := ioutil.ReadFile("./rsa.pub.pem")
	if err != nil {
		panic(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	if err != nil {
		panic(err)
	}
	privateKeyByte, err := ioutil.ReadFile("./rsa.pem")
	if err != nil {
		panic(err)
	}
	privateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, msg, []byte(""))
	if err != nil {
		panic(err)
	}
	rawtext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, []byte(""))
	if err != nil {
		panic(err)
	}
	if string(rawtext) != string(msg) {
		panic(err)
	}
}
