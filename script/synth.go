package script

import (
	"net/http"
	"net/url"
	"bytes"
	"io"
	"github.com/oov/audio"
	"github.com/oov/audio/wave"
)

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
