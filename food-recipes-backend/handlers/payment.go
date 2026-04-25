package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type PaymentRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	TxRef       string `json:"tx_ref"`
	ReturnURL   string `json:"return_url"`
	CallbackURL string `json:"callback_url"`
}

func HandleChapaPayment(w http.ResponseWriter, r *http.Request) {
	var payReq PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&payReq); err != nil {
		http.Error(w, "Invalid payment input", http.StatusBadRequest)
		return
	}

	if payReq.TxRef == "" {
		payReq.TxRef = fmt.Sprintf("txn-%d", time.Now().Unix())
	}

	chapaSecret := os.Getenv("CHAPA_SECRET")
	paymentURL := "https://api.chapa.co/v1/transaction/initialize"

	paymentBody := map[string]string{
		"amount":       payReq.Amount,
		"currency":     payReq.Currency,
		"email":        payReq.Email,
		"first_name":   payReq.FirstName,
		"last_name":    payReq.LastName,
		"tx_ref":       payReq.TxRef,
		"return_url":   payReq.ReturnURL,
		"callback_url": payReq.CallbackURL,
	}

	payload, _ := json.Marshal(paymentBody)

	req, _ := http.NewRequest("POST", paymentURL, strings.NewReader(string(payload)))
	req.Header.Set("Authorization", "Bearer "+chapaSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Chapa request failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
