package main

import (
	"github.com/canack/telebots/services/minipolly/speech"
	"github.com/canack/telebots/services/minipolly/telegram"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var token string

func main() {

	if tokenEnv := os.Getenv("BOT_TOKEN"); tokenEnv == "" {
		panic("Token is not declared.\nPlease attach your token as environment variable. Eg: BOT_TOKEN='token'")
	} else {
		token = tokenEnv
	}

	log.Println("Bot started")
	cleanDirs()
	handleSigterm()
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

// delete all directories if they're exists
func deleteDirs() {
	if _, err := os.Stat("tmp"); err == nil {
		os.RemoveAll("tmp")
	}
	if _, err := os.Stat("videos"); err == nil {
		os.RemoveAll("videos")
	}
}

// create clean directories
func cleanDirs() {
	deleteDirs()
	createDirs()
}

// handle sigterm
func handleSigterm() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Bot stopped")
		cleanDirs()
		os.Exit(1)
	}()
}
