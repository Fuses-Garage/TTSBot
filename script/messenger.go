package script

import (
	"log"
	"github.com/bwmarrin/discordgo"
)
func SendMessage(s *discordgo.Session, channelID string, msg string)(e error) {//メッセージを指定されたチャンネルに投稿します
	_, err := s.ChannelMessageSend(channelID, msg)//送ります
	log.Println(">>> " + msg)
	return err//エラーデータを返します
}
