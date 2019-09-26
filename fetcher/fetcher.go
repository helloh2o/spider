package fetcher

import (
	lg "crawler/log"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var (
	// 5 requests per second is ok
	limiter = time.Tick(time.Millisecond * 150)
	client  http.Client
)

const (
	//Googlebot
	user_agent = "Baiduspider" //"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.75 Safari/537.36"
)

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar
	client.Timeout = time.Second * 10
}

func Fetch(url string) ([]byte, error) {
	<-limiter
	lg.Printf("fetching url %s", url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("user-agent", user_agent)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
