package main

import (
	"fmt"
	"testing"

	"github.com/aliforever/go-pasargad/pasargad"
	"github.com/kr/pretty"
)

func TestPasargad_GeneratePayment(t *testing.T) {
	rsaKey := `YOUR PRIVATE RSA XML HERE`
	p := pasargad.Pasargad{MerchantCode: "4485692", TerminalCode: "1666125", RSAPrivateKey: rsaKey}
	pay, err := p.GeneratePayment(1, 12990, "http://irangopher.com/go/pay", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	pretty.Println(pay)
}
