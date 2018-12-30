package pasargad

import "encoding/xml"

type PaymentResult struct {
	XMLName                xml.Name `xml:"resultObj"`
	Result                 bool     `xml:"result"`
	Action                 string   `xml:"action"`
	TransactionReferenceId int64    `xml:"transactionReferenceID"`
	InvoiceNumber          int      `xml:"invoiceNumber"`
	InvoiceDate            string   `xml:"invoiceDate"`
	MerchantCode           string   `xml:"merchantCode"`
	TerminalCode           string   `xml:"terminalCode"`
	Amount                 int      `xml:"amount"`
}

func (pr *PaymentResult) Unmarshal(xmlBytes []byte) (err error) {
	err = xml.Unmarshal(xmlBytes, &pr)
	return
}
