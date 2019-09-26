package zhenai

import (
	"bytes"
	"crawler/engine"
	"github.com/PuerkitoBio/goquery"
)

func ParseCityList(_ *engine.Request, contents []byte) engine.ParseResult {
	result := engine.ParseResult{}
	dom, err := goquery.NewDocumentFromReader(bytes.NewBuffer(contents))
	if err != nil {
		return result
	}
	nodes := dom.Find(".city-list").Find("a").Nodes
	for _, node := range nodes {
		//<a href="http://www.zhenai.com/zhenghun/ali" data-v-5e16505f>阿里</a
		city := node.FirstChild.Data
		href := node.Attr[0].Val
		if href != "_blank" {
			result.Requests = append(result.Requests, engine.Request{
				Url:       href,
				ParserFun: ParseCity,
			})
			result.Items = append(result.Items, city)
		}
	}
	return result
}
