package main

import (
	"github.com/yanzay/tbot"
	"log"
	"os"
	"telegram_wallet/handlers"
)

func main() {
	bot, err := tbot.NewServer(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	about := ""
	bot.Handle("/about", about)

	bot.HandleFunc("/sign_in", handlers.SignInHandler)
	bot.HandleFunc("/log_in", handlers.LogInHandler)
	bot.HandleFunc("/get_balance", handlers.GetBalanceHandler)
	bot.HandleFunc("/pay_in", handlers.PayInHandler)
	bot.HandleFunc("/pay_out", handlers.PayOutHandler)
	bot.HandleFunc("/transfer", handlers.TransferHandler)

	bot.ListenAndServe()
}
