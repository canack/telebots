package telegram

import (
	"github.com/canack/telebots/services/bigpolly/telegram/handler"
	tele "gopkg.in/telebot.v3"
	_ "image/jpeg"
	_ "image/png"

	"time"
)

var bot *tele.Bot

func SetupTelegramBot(token string) error {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	var err error
	bot, err = tele.NewBot(pref)
	if err != nil {
		return err
	}

	handler.SetupBotHandlers(bot)

	return nil
}

func StartTelegramBot() {
	bot.Start()
}
