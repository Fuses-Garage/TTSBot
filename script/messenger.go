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
func SendEmbedWithField(s *discordgo.Session, channelID,title,desc string,field []*discordgo.MessageEmbedField)(e error) {//埋め込みメッセージを指定されたチャンネルに投稿します
	embed:=&discordgo.MessageEmbed{
		Author:&discordgo.MessageEmbedAuthor{},
		Color:0x880088,
		Title:title,
		Description:desc,
		Fields: field,
	}
	_,err:=s.ChannelMessageSendEmbed(channelID,embed)
	return err//エラーデータを返します
}
func SendEmbed(s *discordgo.Session, channelID,title,desc string)(e error) {//埋め込みメッセージを指定されたチャンネルに投稿します
	embed:=&discordgo.MessageEmbed{
		Author:&discordgo.MessageEmbedAuthor{},
		Color:0x880088,
		Title:title,
		Description:desc,
	}
	_,err:=s.ChannelMessageSendEmbed(channelID,embed)
	return err//エラーデータを返します
}
