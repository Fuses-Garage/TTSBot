package script

import (
	"net/http"
	"net/url"
	"bytes"
	"io"
	"github.com/bwmarrin/discordgo"
)
const fwMS = 20
var vcsession *discordgo.VoiceConnection
func GetBinary(s string)([]byte) {
	url_query:="http://localhost:50021/audio_query?text="+url.QueryEscape(s)+"&speaker=1"
	req, _ := http.NewRequest("POST", url_query, nil)
	req.Header.Set("accept","application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err!=nil {
		panic(err)
	}
	url_synth:="http://localhost:50021/synthesis?speaker=1&enable_interrogative_upspeak=true"
	req_s, _ := http.NewRequest("POST", url_synth, resp.Body)
	req_s.Header.Set("accept","application/json")
	req_s.Header.Set("Content-Type","application/json")
	resp_s, err := client.Do(req_s)
	if err!=nil {
		panic(err)
	}
	defer resp_s.Body.Close()
	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, resp_s.Body); err != nil {
		panic(err)
	}
	return buff.Bytes()
}

func WaveToOpus(b []byte){

}
func Play(opus [][]byte,vc *discordgo.VoiceConnection) {
	vc.Speaking(true)
	defer vc.Speaking(false)
	for _, f := range opus {
		vc.OpusSend <- f
	}
	
}
func Connect(s *discordgo.Session,m *discordgo.MessageCreate){
	userstate,_:=s.State.VoiceState(m.GuildID,m.Author.ID)
	if userstate==nil{
		SendEmbed(s,m.ChannelID, "エラーが発生しました。","呼び出す前にVCに参加してください。")
		return
	}
	var err error=nil
	vcsession, err = s.ChannelVoiceJoin(m.GuildID,userstate.ChannelID, false, false)
	if err!=nil{
		SendEmbed(s,m.ChannelID, "エラーが発生しました。",err.Error())
	}else{
		txtchan,_:=s.Channel(m.ChannelID)
		tcname:=txtchan.Name
		voicechan,_:=s.Channel(userstate.ChannelID)
		vcname:=voicechan.Name
		field:=[]*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "読み上げ元",
				Value:tcname,
			},
			&discordgo.MessageEmbedField{
				Name: "読み上げ先",
				Value:vcname,
			},
		}
		SendEmbedWithField(s,m.ChannelID, "読み上げ開始","これより、読み上げを開始します。",field)
	}
}
func Disconnect(s *discordgo.Session,m *discordgo.MessageCreate){
	
	if vcsession==nil{
		SendEmbed(s,m.ChannelID, "エラーが発生しました。","まだVCに参加していません。")
		return
	}
	err:=vcsession.Disconnect()
	if err!=nil{
		SendEmbed(s,m.ChannelID, "エラーが発生しました。",err.Error())
	}else{
		SendEmbed(s,m.ChannelID, "退出完了","正常に退出しました。")
	}
	
}