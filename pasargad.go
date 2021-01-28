package go_pasargad

import (
	"fmt"
	"time"

	"github.com/aliforever/go-pasargad/objects"

	"net/url"

	"net/http"

	"strings"

	"io/ioutil"

	"github.com/aliforever/go-pasargad/rsa"
)

type Pasargad struct {
	MerchantCode  string
	TerminalCode  string
	RSAPrivateKey string
}

func (p *Pasargad) signData(data string) (result string, err error) {
	xm := rsa.XMLRSAKey{}
	res, err := xm.ConvertToKey(p.RSAPrivateKey)
	if err != nil {
		return
	}
	signer, err := rsa.NewSigner(res)
	if err != nil {
		return
	}
	result, err = signer.SignBase64([]byte(data))
	return
}

func (p *Pasargad) GeneratePayment(invoiceNumber, amount int, redirectUrl string, autoRedirect bool) (paymentForm string, err error) {
	loc, _ := time.LoadLocation("Asia/Tehran")
	timeStr := time.Now().In(loc).Format("2006-01-02 15:04:05")
	data := fmt.Sprintf("#%s#%s#%d#%s#%d#%s#%d#%s#", p.MerchantCode, p.TerminalCode, invoiceNumber, timeStr, amount, redirectUrl, PaymentRequest, timeStr)
	sign, err := p.signData(data)
	if err != nil {
		return
	}

	frm := `<form id="pcp-form2" class="pcp-hidden" method='post' action='https://pep.shaparak.ir/gateway.aspx'>
    <input type='hidden' name='invoiceNumber' value='%d' >
    <input type='hidden' name='invoiceDate' value='%s' >
    <input type='hidden' name='amount' value='%d' >
    <input type='hidden' name='terminalCode' value='%s' >
    <input type='hidden' name='merchantCode' value='%s' >
    <input type='hidden' name='redirectAddress' value='%s' >
    <input type='hidden' name='timeStamp' value='%s' >
    <input type='hidden' name='action' value='%d' >
    <input type='hidden' name='sign' value='%s' >
	</form>`
	paymentForm = fmt.Sprintf(frm, invoiceNumber, timeStr, amount, p.TerminalCode, p.MerchantCode, redirectUrl, timeStr, PaymentRequest, sign)
	if autoRedirect {
		paymentForm += "<script>document.getElementById('pcp-form2').submit();</script>"
	}

	return
}

func (p *Pasargad) GetPaymentResult(invoiceNumber int) (pr *objects.PaymentResult, err error) {
	values := url.Values{}
	values.Add("invoiceUID", fmt.Sprintf("%d", invoiceNumber))
	req, err := http.Post("https://pep.shaparak.ir/CheckTransactionResult.aspx", "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
	if err != nil {
		return
	}
	defer req.Body.Close()
	resp, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	pr = &objects.PaymentResult{}
	err = pr.Unmarshal(resp)
	return
}

func (p *Pasargad) VerifyPayment(amount, invoiceNumber int, timeStr string) (vpr *objects.VerifyPaymentResult, err error) {
	loc, _ := time.LoadLocation("Asia/Tehran")
	timeNow := time.Now().In(loc).Format("2006-01-02 15:04:05")
	values := url.Values{}
	values.Add("MerchantCode", p.MerchantCode)
	values.Add("TerminalCode", p.TerminalCode)
	values.Add("InvoiceNumber", fmt.Sprintf("%d", invoiceNumber))
	values.Add("InvoiceDate", timeStr)
	values.Add("amount", fmt.Sprintf("%d", amount))
	values.Add("TimeStamp", timeNow)

	data := "#%s#%s#%d#%s#%d#%s#"
	data = fmt.Sprintf(data, p.MerchantCode, p.TerminalCode, invoiceNumber, timeStr, amount, timeNow)
	sign, err := p.signData(data)
	if err != nil {
		return
	}
	values.Add("sign", sign)
	req, err := http.Post("https://pep.shaparak.ir/VerifyPayment.aspx", "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
	if err != nil {
		return
	}
	defer req.Body.Close()
	res, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	vpr = &objects.VerifyPaymentResult{}
	err = vpr.Unmarshal(res)
	return
}
