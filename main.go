package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	arif "github.com/AnaniyaBelew/ArifpayGoPlugin"
	cbe "github.com/yishak-cs/BravoArif/Deposit/CBE"
	telebirr "github.com/yishak-cs/BravoArif/Deposit/telebirr"
)

type Response struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
	Data  Data   `json:"data"`
}

type Data struct {
	SessionID   string  `json:"sessionId"`
	PaymentURL  string  `json:"paymentUrl"`
	CancelURL   string  `json:"cancelUrl"`
	TotalAmount float64 `json:"totalAmount"`
}

func generateNonce() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("unable to generate nonce")
	}
	return hex.EncodeToString(bytes), nil
}

func testTelebirrDeposit() {
	nonce, err := generateNonce()
	if err != nil {
		fmt.Println("unable to generate nonce", err)
		os.Exit(0)
	}
	paymentReq := arif.PaymentRequest{
		Phone:          "0907968056",
		Nonce:          nonce,
		Email:          "telebirrTest@gmail.com",
		CancelUrl:      "http://example.com",
		ErrorUrl:       "http://error.com",
		SuccessUrl:     "http://example.com",
		NotifyUrl:      "https://gateway.arifpay.net/test/callback",
		PaymentMethods: []string{"TELEBIRR_USSD"},
		ExpireDate:     time.Now().AddDate(0, 0, 1),
		Items: []interface{}{
			map[string]interface{}{
				"name":     "bet",
				"quantity": 2,
				"price":    20,
			},
		},
		Beneficiaries: []struct {
			AccountNumber string  `json:"accountNumber"`
			Bank          string  `json:"bank"`
			Amount        float64 `json:"amount"`
		}{
			{
				AccountNumber: "01320811436100",
				Bank:          "AWINETAA",
				Amount:        40.0,
			},
		},
		Lang: "EN",
	}
	ResponseStr, err := telebirr.Deposit(&paymentReq)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(ResponseStr)
}

func main() {
	nonce, err := generateNonce()
	if err != nil {
		fmt.Println("unable to generate nonce", err)
		os.Exit(0)
	}
	paymentReq := arif.PaymentRequest{
		Phone:          "251912345678",
		Email:          "Test@gmail.com",
		Nonce:          nonce,
		CancelUrl:      "https://example.com",
		ErrorUrl:       "http://error.com",
		SuccessUrl:     "http://example.com",
		NotifyUrl:      "https://example.com",
		PaymentMethods: []string{"CBE"},
		ExpireDate:     time.Date(2025, 2, 39, 29, 45, 27, 0, time.UTC),
		Items: []interface{}{
			map[string]interface{}{
				"name":        "Bet",
				"quantity":    1,
				"price":       20,
				"description": "bet on BravoBet",
			},
		},
		Beneficiaries: []struct {
			AccountNumber string  `json:"accountNumber"`
			Bank          string  `json:"bank"`
			Amount        float64 `json:"amount"`
		}{
			{
				AccountNumber: "01320811436100",
				Bank:          "AWINETAA",
				Amount:        20.0,
			},
		},
		Lang: "EN",
	}

	responseStr, err := cbe.Deposit(&paymentReq)
	if err != nil {
		fmt.Printf("Error making deposit request: %v\n", err)
		return
	}

	if responseStr == "" {
		fmt.Println("Received empty response from API")
		return
	}

	fmt.Printf("Raw response: %s\n", responseStr)

	var response Response
	prob := json.Unmarshal([]byte(responseStr), &response)
	if prob != nil {
		fmt.Printf("Error parsing JSON: %v\n", prob)
		return
	}

	// Extract the session ID
	sessionID := response.Data.SessionID
	fmt.Println("\nSession ID:", sessionID)

	cbeDep := cbe.CBERequest{
		SessionID:   sessionID,
		PhoneNumber: "251912345678",
		Password:    "cbe123",
	}

	actual, err := cbe.SecondReq(&cbeDep)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(actual)
}
