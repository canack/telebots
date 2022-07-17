package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/canack/telebots/services/bigpolly/page"
	tele "gopkg.in/telebot.v3"
	"image/jpeg"
	"regexp"
)

func checkSite(site string) bool {

	redditExp, err := regexp.Compile(`^(https?)://(www)?.?reddit.com.*`)
	if err != nil {
		panic(err)
	}
	if redditExp.MatchString(site) {
		return true
	}
	return false

}

func Crawl(c tele.Context) error {
	if !checkSite(c.Text()) {
		return c.Send("This site not supported for now.")
	}

	img, err := page.CrawlReddit(c.Text())
	if errors.Is(err, context.DeadlineExceeded) {
		return c.Send("Process time out.")
	} else if err != nil {
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
