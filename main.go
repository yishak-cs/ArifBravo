package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	arif "github.com/AnaniyaBelew/ArifpayGoPlugin"
	deposit "github.com/yishak-cs/BravoArif/Deposit/CBE"
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

func main() {
	nonce, err := generateNonce()
	if err != nil {
		fmt.Println("unable to generate nonce", err)
		os.Exit(0)
	}
	paymentReq := arif.PaymentRequest{
		Phone:          "0907968056",
		Nonce:          nonce,
		CancelUrl:      "http://example.com",
		ErrorUrl:       "http://error.com",
		SuccessUrl:     "http://example.com",
		NotifyUrl:      "https://gateway.arifpay.net/test/callback",
		PaymentMethods: []string{"CBE"},
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

	responseStr, _ := deposit.Deposit(&paymentReq)
	fmt.Println(responseStr)
	var response Response
	prob := json.Unmarshal([]byte(responseStr), &response)
	if prob != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Extract the session ID
	sessionID := response.Data.SessionID
	fmt.Println("\nSession ID:", sessionID)

	cbe := deposit.CBERequest{
		SessionID:   sessionID,
		PhoneNumber: "0907968056",
		Password:    "cbe123",
	}

	actual, err := deposit.SecondReq(&cbe)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(actual)
}
