package cbe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	arif "github.com/AnaniyaBelew/ArifpayGoPlugin"
)

const apiKey = "my-api-key"

func Deposit(req *arif.PaymentRequest) (string, error) {
	payment := arif.NewPayment(apiKey, req.ExpireDate)

	paymentRequestBytes, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	fmt.Printf("Request payload: %s\n", string(paymentRequestBytes))

	httpreq, err := http.NewRequest("POST", "https://gateway.arifpay.net/api/checkout/session", bytes.NewBuffer(paymentRequestBytes))
	if err != nil {
		return "", fmt.Errorf("create request error: %w", err)
	}

	httpreq.Header.Set("Content-Type", "application/json")
	httpreq.Header.Set("x-arifpay-key", payment.APIKey)

	fmt.Printf("Request headers: %+v\n", httpreq.Header)

	client := &http.Client{}
	resp, err := client.Do(httpreq)
	if err != nil {
		return "", fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response error: %w", err)
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response body: %s\n", string(responseBytes))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status code: %d, body: %s", resp.StatusCode, string(responseBytes))
	}

	return string(responseBytes), nil
}
