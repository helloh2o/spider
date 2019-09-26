package zhenai

import (
	"bytes"
	"crawler/engine"
	lg "crawler/log"
	"github.com/PuerkitoBio/goquery"
)

func ParseCity(req *engine.Request, contents []byte) engine.ParseResult {
	defer func() {
		if r := recover(); r != nil {
			lg.Printf("recover from ParseCity %v", r)
			//lg.Printf("Content \n %s", string(contents))
		}
	}()
	result := engine.ParseResult{}
	dom, err := goquery.NewDocumentFromReader(bytes.NewBuffer(contents))
	if err != nil {
		return result
	}
	contentNode := dom.Find(".list-item").Find(".content")
	nodes := contentNode.Find("a").Nodes
	genders := contentNode.Find("td:contains('性别')").Nodes
	for i, node := range nodes {
		//<a href="http://album.zhenai.com/u/1810062308" target="_blank">谭佑霆</a>
		name := node.FirstChild.Data
		href := node.Attr[0].Val
		if href != "_blank" {
			gender := genders[i].LastChild.Data
			result.Requests = append(result.Requests, engine.Request{
				Url:       href,
				ParserFun: ParseProfile,
				Attach:    gender,
			})
			result.Items = append(result.Items, name)
		}
	}
	// next page
	pages := dom.Find(".paging-item a").Nodes
	for _, p := range pages {
		if len(p.Attr) == 0 {
			continue
		}
		href := p.Attr[0].Val
		//lg.Printf("Add next page %s", href)
		result.Requests = append(result.Requests, engine.Request{
			Url:       href,
			ParserFun: ParseCity,
		})
	}
	// other items
	items := dom.Find(".city-list .list-item a:contains('征婚')").Nodes
	//lg.Printf("items %d", len(items))
	for _, item := range items {
		if len(item.Attr) < 2 {
			continue
		}
		href := item.Attr[1].Val
		//lg.Printf("Add item url %s", href)
		result.Requests = append(result.Requests, engine.Request{
			Url:       href,
			ParserFun: ParseCity,
		})
	}
	return result
}
