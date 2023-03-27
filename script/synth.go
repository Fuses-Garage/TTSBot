package script

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode/utf8"
)

type vcData struct { //読み上げデータ
	connection *discordgo.VoiceConnection //音声を再生するコネクション
	channelID  string                     //読み上げるテキストチャンネルのID
	queue      *[]string                  //読み上げたい音声のパスのキュー（のアドレス）
}

var vcDict map[string]vcData //サーバごとの読み上げデータ
func TTS(m *discordgo.MessageCreate) { //渡されたメッセージを読み上げる
	go func() { //ゴルーチンを使う
		v, ok := vcDict[m.GuildID] //接続データを取得
		if ok {                    //VC接続中なら
			if v.channelID == m.ChannelID { //書き込まれた先が読み上げ対象なら
				b, err := getBinary(m.Content) //バイナリデータを取得
				if err != nil {                //エラーが起きたら
					log.Printf("error: %v", err)
					return
				}
				var path string                        //ファイルのパス
				path, err = makeWaveFile(b, m.GuildID) //ファイルに書き込む
				if err != nil {                        //エラーが起きたら
					log.Printf("error: %v", err)
					return
				}
				play(&v, path, false) //再生
			}
		}
	}()
}
func getBinary(s string) ([]byte, error) { //バイナリデータをもらってくる
	str := s
	if utf8.RuneCountInString(str) > 30 { //あまりにも文章が長いときは
		slice := []rune(str)                            //ルーンに変換しないと二バイト文字がバグる
		strarr := []string{string(slice[:30]), "いかりゃく"} //30文字で切って以下略をつける
		str = strings.Join(strarr, "　")                 //くっ付ける
	}
	urlParts := []string{"http://localhost:50021/audio_query?text=", url.QueryEscape(str), "&speaker=8"}
	url_query := strings.Join(urlParts, "")           //URL組み立て
	req, _ := http.NewRequest("POST", url_query, nil) //POSTでリクエスト
	req.Header.Set("accept", "application/json")      //ヘッダをセット

	client := new(http.Client)  //クライアント生成
	resp, err := client.Do(req) //リクエスト
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	url_synth := "http://localhost:50021/synthesis?speaker=8&enable_interrogative_upspeak=true" //音声生成用URL
	req_s, _ := http.NewRequest("POST", url_synth, resp.Body)                                   //POSTでリクエスト
	req_s.Header.Set("accept", "audio/wav")                                                     //ヘッダをセット
	req_s.Header.Set("Content-Type", "application/json")                                        //ヘッダをセット
	resp_s, err := client.Do(req_s)                                                             //リクエスト
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	defer resp_s.Body.Close()                             //いらなくなったら閉じる
	buff := bytes.NewBuffer(nil)                          //バッファ生成
	if _, err := io.Copy(buff, resp_s.Body); err != nil { //移す
		log.Printf("error: %v", err)
		return nil, err
	}
	return buff.Bytes(), nil //バッファを返す
}

func makeWaveFile(b []byte, guildID string) (string, error) { //音声を生成する関数
	max := new(big.Int)                  //ファイル名の重複を避けるための乱数
	max.SetInt64(int64(1000000))         //100万通り
	r, err := rand.Int(rand.Reader, max) //乱数生成
	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}
	path := fmt.Sprintf("%s_%d.wav", guildID, r) //ファイル名の重複を避ける
	file, _ := os.Create(path)                   //ファイル生成
	defer func() {
		file.Close() //終わったら閉じる
	}()
	file.Write(b)    //ファイルにデータを書き込む
	return path, nil //ファイルのパスを返す
}
func play(v *vcData, path string, force bool) { //指定されたパスの音声を再生またはキューに追加します。
	if len(*v.queue) > 0 && !force { //現在再生中かつ強制再生フラグがオフなら
		*v.queue = append(*v.queue, path) //キューに追加
		return
	}
	*v.queue = append(*v.queue, path) //キューに追加
	vc := v.connection                //コネクションのデータを取り出す
	vc.Speaking(true)                 //再生開始
	defer vc.Speaking(false)          //終わったら再生フラグを戻す
	defer os.Remove(path)             //終わったらファイルを消す
	defer func(v *vcData) {           //終わったら
		*v.queue = (*v.queue)[1:] //キューから終わった分を消す
		if len(*v.queue) > 0 {    //まだキューに残ってたら
			pt := (*v.queue)[0] //次に読み上げるパスを取得
			*v.queue = (*v.queue)[1:]
			play(v, pt, true) //強制再生
		}
	}(v)
	dgvoice.PlayAudioFile(vc, path, make(chan bool)) //再生開始
}
func Connect(s *discordgo.Session, m *discordgo.MessageCreate) {
	if vcDict == nil { //vcDictがnilなら
		vcDict = make(map[string]vcData) //連想配列を生成
	}

	userstate, _ := s.State.VoiceState(m.GuildID, m.Author.ID) //呼び出したユーザはVCに入っているか？
	if userstate == nil {                                      //入っているなら
		SendEmbed(s, m.ChannelID, "エラーが発生しました。", "呼び出す前にVCに参加してください。")
		return
	}
	_, ok := vcDict[m.GuildID] //そのサーバーの接続データは存在するか
	if ok {                    //存在するなら
		SendEmbed(s, m.ChannelID, "エラーが発生しました。", "すでにVCに接続中です。")
		return
	}
	vcsession, err := s.ChannelVoiceJoin(m.GuildID, userstate.ChannelID, false, false) //VCに接続
	if err != nil {
		SendEmbed(s, m.ChannelID, "エラーが発生しました。", err.Error())
		return
	} else {
		txtchan, _ := s.Channel(m.ChannelID)
		tcname := txtchan.Name
		voicechan, _ := s.Channel(userstate.ChannelID)
		vcname := voicechan.Name
		field := []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "読み上げ元",
				Value: tcname,
			},
			&discordgo.MessageEmbedField{
				Name:  "読み上げ先",
				Value: vcname,
			},
		}
		SendEmbedWithField(s, m.ChannelID, "読み上げ開始", "これより、読み上げを開始します。", field)
		slice := make([]string, 0, 10)
		newData := vcData{vcsession, m.ChannelID, &slice} //接続データを生成
		vcDict[m.GuildID] = newData                       //サーバIDに対応付ける
	}
}
func Disconnect(s *discordgo.Session, m *discordgo.MessageCreate) { //VCから抜ける
	v, ok := vcDict[m.GuildID] //データ取得
	if !ok {                   //データが存在しない=VCに入ってない
		SendEmbed(s, m.ChannelID, "エラーが発生しました。", "まだVCに参加していません。")
		return
	}
	err := v.connection.Disconnect() //退出
	if err != nil {                  //退出時にエラーが起きたら
		SendEmbed(s, m.ChannelID, "エラーが発生しました。", err.Error())
	} else { //起きなかったら
		delete(vcDict, m.GuildID) //接続データを削除
		SendEmbed(s, m.ChannelID, "退出完了", "正常に退出しました。")
	}
}
