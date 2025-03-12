package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/yanzay/tbot/v2"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"telegram_wallet_bot/internal"
)

type PaymentBody struct {
	Amount  float64 `json:"amount"`
	OrderId string  `json:"order_id"`
}

type PaymentResponse struct {
	How     string `json:"how"`
	Request string `json:"request"`
}

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")

	merchantID := os.Getenv("MERCHANT_ID")
	pubKeyID := os.Getenv("PUB_KEY_ID")
	priKeyPath := os.Getenv("PRI_KEY_PATH")
	baseUrl := os.Getenv("BASE_URL")
	pathDeposit := os.Getenv("PATH_DEPOSIT")

	bot := tbot.New(token)
	c := bot.Client()
	bot.HandleMessage("/payin .+", func(m *tbot.Message) {
		text := strings.TrimPrefix(m.Text, "/payin ")
		m, err := c.SendMessage(m.Chat.ID, payin(m, merchantID, pubKeyID, priKeyPath, text, baseUrl, pathDeposit))
		if err != nil {
			fmt.Printf(err.Error())
		}
		fmt.Printf(m.Text)
	})

	bot.HandleMessage("/payout .+", func(m *tbot.Message) {
		text := strings.TrimPrefix(m.Text, "/payout ")
		c.SendMessage(m.Chat.ID, payout(m, merchantID, pubKeyID, priKeyPath, text, baseUrl, pathDeposit), tbot.OptParseModeMarkdown)
	})
	bot.Start()
}

func payin(m *tbot.Message, merchantID, pubKeyID, priKeyPath, data, baseUrl, pathDeposit string) string {
	requestURL := baseUrl + pathDeposit
	amount, err := strconv.ParseFloat(strings.TrimSpace(data), 64)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return ""
	}

	payment, err := json.Marshal(PaymentBody{Amount: amount, OrderId: uuid.NewString()})
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return ""
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader([]byte(payment)))
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return ""
	}

	token := internal.GenerateAuthJWT(priKeyPath, merchantID, pubKeyID,
		m.From.FirstName+"_"+
			m.From.LastName+"_"+
			strconv.Itoa(m.From.ID)+
			".telegram")

	req.Header.Add("Authorization", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return ""
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return ""
	}

	var response PaymentResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return ""
	}

	payinLink := fmt.Sprintf("%s?token=%s", response.How, token)
	return payinLink
}

func payout(m *tbot.Message, merchantID, pubKeyID, priKeyPath, data, baseUrl, pathWithdraw string) string {
	requestURL := baseUrl + pathWithdraw
	req, err := http.NewRequest(http.MethodPost, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return ""
	}

	token := internal.GenerateAuthJWT(priKeyPath, merchantID, pubKeyID,
		m.From.FirstName+"_"+
			m.From.FirstName+"_"+
			strconv.Itoa(m.From.ID)+
			".telegram")

	req.Header.Add("Authorization", token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return ""
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("client: error making http request: %s\n", res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		return ""
	}

	return string(resBody)
}
