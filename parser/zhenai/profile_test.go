package zhenai

import (
	"bytes"
	"crawler/engine"
	lg "crawler/log"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestParseProfile(t *testing.T) {
	f, err := os.Open("profile.txt")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, f))
	if err != nil {
		panic(err)
	}
	profile := ParseProfile(&engine.Request{}, data)
	for _, p := range profile.Items {
		lg.Printf("Got item %+v", p)
	}
}

func TestParseProfile2(t *testing.T) {
	resp, err := http.Get("https://album.zhenai.com/u/1027979271")
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	lg.Printf(string(data))
	profile := ParseProfile(&engine.Request{}, data)
	for _, p := range profile.Items {
		lg.Printf("Got item %+v", p)
	}
}

func TestGoQuery(t *testing.T) {
	f, err := os.Open("profile.txt")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, f))
	if err != nil {
		panic(err)
	}
	dom, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	lg.Printf("%d", dom.Find(".m-userInfo").Find(".logo").Size())
}
