package handler

import (
	tele "gopkg.in/telebot.v3"
	"log"
)

type mediaFile struct {
	pitch                                    string
	text                                     string
	rawSoundFilename, pitchedSoundFilename   string
	rawVideoFilename, processedVideoFilename string
	cuttedVideoFilename                      string
	//soundLength as seconds
	soundLength float64
}

func SetupBotHandlers(bot *tele.Bot) {
	bot.Handle("/start", Welcome)
	bot.Handle(tele.OnText, ProcessText)

}

func Welcome(c tele.Context) error {
	c.Notify(tele.Typing)
	return c.Send(`Hi!
	
	Github: github.com/canack`, &tele.SendOptions{DisableWebPagePreview: true})
}

func ProcessText(c tele.Context) error {

	c.Notify(tele.Typing)

	var media mediaFile
	defer media.DeleteTempFiles()
	media.text = c.Text()
	media.pitch = "150"
	media.rawSoundFilename = "tmp/" + randomString(32) + ".mp3"
	media.pitchedSoundFilename = "tmp/" + randomString(32) + ".mp3"
	media.processedVideoFilename = "tmp/" + randomString(32) + ".mp4"

	if err := media.GenerateMedia(); err != nil {
		return c.Send("Eğer gönderdiğin text oldukça kısaysa bu nedenden hata almışsındır. Yine de bu durumu bildir")
	}
	log.Println("Media generated")

	c.Notify(tele.UploadingVideo)
	file := &tele.Video{File: tele.FromDisk(media.processedVideoFilename), FileName: "video.mp4"}
	return c.Send(file)

}
