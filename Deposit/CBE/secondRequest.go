package cbe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CBERequest struct {
	SessionID   string `json:"sessionId"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

func SecondReq(cbe *CBERequest) (string, error) {

	reqBody, err := json.Marshal(cbe)

	fmt.Println(string(reqBody))
	if err != nil {
		return "", err
	}

	httpreq, err := http.NewRequest("POST", "https://gateway.arifpay.net/api/checkout/cbe/direct/transfer", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	httpreq.Header.Set("Content-Type", "application/json")
	httpreq.Header.Set("x-arifpay-key", apiKey)

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
