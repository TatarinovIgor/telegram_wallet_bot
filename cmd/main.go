package main

import (
	"fmt"
	"github.com/yanzay/tbot"
	"log"
	"os"
	"telegram_wallet_bot/handlers"
)

func main() {
	tbBotKey := os.Getenv("TELEGRAM_TOKEN")
	merchantID := os.Getenv("MERCHANT_ID")
	pubKeyID := os.Getenv("PUB_KEY_ID")
	priKeyPath := os.Getenv("PRI_KEY_PATH")

	bot, err := tbot.NewServer(tbBotKey)
	if err != nil {
		log.Fatal(err)
	}

	about := "hi"
	bot.Handle("/about", about)

	bot.HandleFunc("/pay_in", handlers.PayInHandler(merchantID, pubKeyID, priKeyPath))
	bot.HandleFunc("/pay_out", handlers.PayOutHandler(merchantID, pubKeyID, priKeyPath))

	fmt.Printf("Bot initialized, serving at following API key: %s, for %s merchant", tbBotKey, merchantID)
	err = bot.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}
