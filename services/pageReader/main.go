package main

import (
	"github.com/canack/telebots/services/pageReader/telegram"
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

	cleanTmpDir()
	handleSigterm()
	startBot()

}

func startBot() {
	if err := telegram.SetupTelegramBot(token); err != nil {
		panic(err)
	}

	telegram.StartTelegramBot()
}

// handle sigterm
func handleSigterm() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Bot stopped")
		deleteTmpDir()
		os.Exit(1)
	}()
}

// make tmp directory if it's not exists
func makeTmpDir() {
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		os.Mkdir("tmp", 760)
	}
}

// delete tmp directory if it's exists
func deleteTmpDir() {
	if _, err := os.Stat("tmp"); err == nil {
		os.RemoveAll("tmp")
	}
}

// create clean tmp directory
func cleanTmpDir() {
	deleteTmpDir()
	makeTmpDir()
}
