package handler

import (
	"bytes"
	"github.com/canack/telebots/services/page_reader/page"
	tele "gopkg.in/telebot.v3"
	"image/jpeg"
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

func Crawl(c tele.Context) error {
	img, err := page.CrawlReddit(c.Text())
	if err != nil {
		return c.Send(err.Error())
	}

	for _, i := range *img {

		writer := new(bytes.Buffer)

		jpeg.Encode(writer, i.Image, nil)

		photo := tele.Photo{
			File:    tele.FromReader(writer),
			Caption: i.Text,
		}

		c.Reply(&photo)

	}

	return nil

}
