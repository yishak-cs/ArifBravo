package cbe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	arif "github.com/AnaniyaBelew/ArifpayGoPlugin"
)

const apiKey = ""

var expireDate = time.Now().AddDate(0, 0, 1)

func Deposit(req *arif.PaymentRequest) (string, error) {

	payment := arif.NewPayment(apiKey, expireDate)

	fmt.Println("Request sent:", payment)

	paymentRequestBytes, err := json.Marshal(&req)
	if err != nil {
		return "", err
	}
	req.Nonce = fmt.Sprintf("%d", time.Now().UnixNano())
	req.ExpireDate = payment.ExpireDate
	httpreq, err := http.NewRequest("POST", "https://gateway.arifpay.net/api/checkout/session", bytes.NewBuffer(paymentRequestBytes))
	if err != nil {
		return "", err
	}

	httpreq.Header.Set("Content-Type", "application/json")
	httpreq.Header.Set("x-arifpay-key", payment.APIKey)

	client := &http.Client{}
	resp, err := client.Do(httpreq)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBytes), nil
}
