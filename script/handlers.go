package script

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)
var(
	prefix string="!tts"
)
func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {//メッセージが投稿されたら呼ばれます
	u := m.Author//uはmの発信者
	if !u.Bot {//発信元が人間なら
		if(strings.HasPrefix(m.Content,prefix)){
			command:=strings.Split(m.Content," ")
			switch command[1]{
			
				case "s":
					Connect(s,m)//接続処理
				case "e":
					Disconnect(s,m) //今いる通話チャンネルから抜ける
			}
		}else{
			TTS(m)//読み上げ機能を呼び出す
		}
	}
 }
