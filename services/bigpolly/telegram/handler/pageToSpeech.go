package handler

import (
	"context"
	"errors"
	"github.com/canack/telebots/services/bigpolly/mediaProcessing"
	"github.com/canack/telebots/services/bigpolly/page"
	tele "gopkg.in/telebot.v3"
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

	resp, err := page.CrawlReddit(c.Text())
	if errors.Is(err, context.DeadlineExceeded) {
		return c.Send("Process time out.")
	} else if err != nil {
		return c.Send(err.Error())
	}

	final, err := mediaProcessing.Process(resp)
	if err != nil {
		return c.Send(err.Error())
	}

	file := tele.Video{File: tele.FromDisk(final)}
	return c.Send(&file)

}
