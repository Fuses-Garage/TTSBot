package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/Fuses-Garage/TTSBot/script"
	"github.com/bwmarrin/discordgo"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("TTS_BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(script.OnMessageCreate)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// 直近の関数（main）の最後に実行される
	defer dg.Close()

	fmt.Println("Listening...")
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot
}
