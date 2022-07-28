package main

import (
	"github.com/canack/telebots/services/bigpolly/speech"
	"github.com/canack/telebots/services/bigpolly/telegram"
	"github.com/canack/telebots/services/bigpolly/types"
	"io/ioutil"
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
	if err := speech.SetupAWS(); err != nil {
		panic(err)
	}
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
	if _, err := os.Stat(types.TempPath); os.IsNotExist(err) {
		os.Mkdir(types.TempPath, 0760)
	}
}

// delete tmp directory if it's exists
func deleteTmpDir() {
	if _, err := os.Stat(types.TempPath); err == nil {
		os.RemoveAll(types.TempPath)
	}
}

// create clean tmp directory
func cleanTmpDir() {
	deleteTmpDir()
	makeTmpDir()
}

func listFileCountInVideoDir() (int, error) {
	if _, err := os.Stat(types.VideoPath); os.IsNotExist(err) {
		return 0, err
	}
	files, _ := ioutil.ReadDir(types.VideoPath)
	return len(files), nil
}
