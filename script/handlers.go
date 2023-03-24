package script

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {//メッセージが投稿されたら呼ばれます
	u := m.Author//uはmの発信者
	if !u.Bot {//発信元が人間なら
		switch {
			case strings.HasPrefix(m.Content,"!tts s"):
				Connect(s,m)
			case strings.HasPrefix(m.Content,"!tts e"):
				Disconnect() //今いる通話チャンネルから抜ける
		}
	}
 }
