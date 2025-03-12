package handlers

import (
	"github.com/yanzay/tbot"
	"strconv"
	"telegram_wallet_bot/internal"
)

func PayInHandler(merchantID, pubKeyID, priKeyPath string) tbot.HandlerFunction {
	return func(m *tbot.Message) {
		token := internal.GenerateAuthJWT(priKeyPath, merchantID, pubKeyID,
			strconv.Itoa(m.From.ID)+
				".telegram")
		m.Reply(token)

	}
}
