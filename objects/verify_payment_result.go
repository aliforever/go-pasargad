package objects

import "encoding/xml"

type VerifyPaymentResult struct {
	XMLName       xml.Name `xml:"actionResult"`
	Result        bool     `xml:"result"`
	ResultMessage string   `xml:"resultMessage"`
}

func (vpr *VerifyPaymentResult) Unmarshal(xmlBytes []byte) (err error) {
	err = xml.Unmarshal(xmlBytes, &vpr)
	return
}
