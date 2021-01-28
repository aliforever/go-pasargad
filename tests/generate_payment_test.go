package tests

import (
	"fmt"
	"testing"

	go_pasargad "github.com/aliforever/go-pasargad"
)

func TestPasargad_GeneratePayment(t *testing.T) {
	rsaKey := `YOUR PRIVATE RSA XML HERE`
	p := go_pasargad.Pasargad{MerchantCode: "4485692", TerminalCode: "1666125", RSAPrivateKey: rsaKey}
	pay, err := p.GeneratePayment(1, 12990, "http://localhost/go/pay", true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(pay)
}
