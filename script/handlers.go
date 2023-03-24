package script

import (
	"github.com/bwmarrin/discordgo"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {//メッセージが投稿されたら呼ばれます
	u := m.Author//uはmの発信者
	if !u.Bot {//発信元が人間なら
		SendMessage(s, m.ChannelID, m.Content)//メッセージをオウム返し
	}
 
 }
