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
	dg, err := discordgo.New("Bot " + os.Getenv("TTS_BOT_TOKEN"))//トークンは環境変数にしてログイン
	if err != nil {//エラー発生時
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(script.OnMessageCreate)//メッセージが投稿されたら実行
	err = dg.Open()//セッション開始
	if err != nil {//エラーが起きたら
		fmt.Println("error opening connection,", err)
		return
	}
	// 終了時にちゃんとクローズするように
	defer dg.Close()

	fmt.Println("ログイン成功!")
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot
}
