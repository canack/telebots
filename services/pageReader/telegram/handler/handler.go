package handler

import (
	tele "gopkg.in/telebot.v3"
)

func SetupBotHandlers(bot *tele.Bot) {
	bot.Handle("/start", Welcome)
	bot.Handle(tele.OnText, Crawl)
}

func Welcome(c tele.Context) error {
	c.Notify(tele.Typing)
	return c.Send(`Hi!
	
	Github: github.com/canack`, &tele.SendOptions{DisableWebPagePreview: true})
}
