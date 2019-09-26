package zhenai

import (
	"crawler/engine"
	lg "crawler/log"
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseCityList(t *testing.T) {
	const expectedSize = 470
	f, err := os.Open("city.html")
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(f)
	list := ParseCityList(&engine.Request{}, bytes)
	lg.Printf("city size %d", len(list.Items))
	if len(list.Requests) != expectedSize {
		panic(errors.New("not expected size 470"))
	}
}

func TestParseCity(t *testing.T) {
	const expectedSize = 20
	f, err := os.Open("page.html")
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(f)
	list := ParseCity(nil, bytes)
	lg.Printf("user size %d", len(list.Items))
	if len(list.Items) != expectedSize {
		panic(errors.New("not expected size 20"))
	} else {
		for i, name := range list.Items {
			lg.Printf("Got name %s, profile url %s", name, list.Requests[i].Url)
		}
	}
}
