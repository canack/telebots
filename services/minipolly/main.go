package main

import (
	"github.com/canack/telebots/services/minipolly/speech"
	"github.com/canack/telebots/services/minipolly/telegram"
	"log"
	"os"
)

var token string

func main() {

	createDirs()

	if tokenEnv := os.Getenv("BOT_TOKEN"); tokenEnv == "" {
		panic("Token is not declared.\nPlease attach your token as environment variable. Eg: BOT_TOKEN='token'")
	} else {
		token = tokenEnv
	}

	log.Println("Bot started")

	startBot()

}

func startBot() {
	if err := telegram.SetupTelegramBot(token); err != nil {
		panic(err)
	}

	if err := speech.SetupAWS(); err != nil {
		panic(err)
	}

	telegram.StartTelegramBot()
}

// create tmp directory if it's not exists
func createTmpDir() {
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		os.Mkdir("tmp", 0755)
	}
}

// create videos directory if it's not exists
func createVideosDir() {
	if _, err := os.Stat("videos"); os.IsNotExist(err) {
		os.Mkdir("videos", 0755)
	}
}

// create videos/raw directory if it's not exists
func createRawVideosDir() {
	if _, err := os.Stat("videos/raw"); os.IsNotExist(err) {
		os.Mkdir("videos/raw", 0755)
	}
}

// create all directories if they're not exists
func createDirs() {
	createTmpDir()
	createVideosDir()
	createRawVideosDir()
}

//