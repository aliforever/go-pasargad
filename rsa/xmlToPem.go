package rsa

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"math/big"

	"github.com/go-errors/errors"
)

type XMLRSAKey struct {
	Modulus  string
	Exponent string
	P        string
	Q        string
	DP       string
	DQ       string
	InverseQ string
	D        string
	Key      []byte
}

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (x *XMLRSAKey) b64d(str string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(str)
	chkErr(err)
	return []byte(decoded)
}

func (x *XMLRSAKey) b64bigint(str string) *big.Int {
	bInt := &big.Int{}
	bInt.SetBytes(x.b64d(str))
	return bInt
}

func (x *XMLRSAKey) ToPrivateKey(rsaPrivateByte []byte) (key *rsa.PrivateKey, err error) {
	p, _ := pem.Decode(rsaPrivateByte)
	if p == nil {
		return nil, errors.New("no pem block found")
	}
	res, err := x509.ParsePKCS8PrivateKey(p.Bytes)
	if err != nil {
		return
	}
	key = res.(*rsa.PrivateKey)
	return
}

func (x *XMLRSAKey) ToSHA1(data []byte) (res []byte, err error) {
	h := sha1.New()
	_, err = h.Write(data)
	if err != nil {
		return
	}
	res = h.Sum(nil)
	return
	//privateKey10 := privateKey10I.(*rsa.PrivateKey)
	//sig, _ := rsa.SignPKCS1v15(rand.Reader, privateKey10, crypto.SHA1, sum)
	//code := hex.EncodeToString(sig)
	//str = b64.StdEncoding.EncodeToString([]byte(code))
	//return
}

func (x *XMLRSAKey) ConvertToKey(xmlBs string) (block []byte, err error) {
	/*xmlbs, err := ioutil.ReadAll(os.Stdin)
	chkErr(err)

	if decoded, err := base64.StdEncoding.DecodeString(string(xmlbs)); err == nil {
		xmlbs = decoded
	}*/
	xrk := XMLRSAKey{}
	err = xml.Unmarshal([]byte(xmlBs), &xrk)
	if err != nil {
		return
	}
	key := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: xrk.b64bigint(xrk.Modulus),
			E: int(xrk.b64bigint(xrk.Exponent).Int64()),
		},
		D:      xrk.b64bigint(xrk.D),
		Primes: []*big.Int{xrk.b64bigint(xrk.P), xrk.b64bigint(xrk.Q)},
		Precomputed: rsa.PrecomputedValues{
			Dp:        xrk.b64bigint(xrk.DP),
			Dq:        xrk.b64bigint(xrk.DQ),
			Qinv:      xrk.b64bigint(xrk.InverseQ),
			CRTValues: ([]rsa.CRTValue)(nil),
		},
	}

	pemBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	block = pem.EncodeToMemory(pemBlock)
	x.Key = block
	return
}
