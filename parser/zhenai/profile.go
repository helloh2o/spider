package zhenai

import (
	"bytes"
	"crawler/engine"
	lg "crawler/log"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

func ParseProfile(r *engine.Request, contents []byte) engine.ParseResult {
	defer func() {
		if r := recover(); r != nil {
			lg.Printf("recover from ParseProfile %v", r)
			//lg.Printf("Content \n %s", string(contents))
		}
	}()
	result := engine.ParseResult{}
	dom, err := goquery.NewDocumentFromReader(bytes.NewBuffer(contents))
	if err != nil {
		return result
	}
	p := Profile{}
	p.Photo, _ = dom.Find(".m-userInfo .logo").First().Attr("style")
	des := strings.Replace(dom.Find(".m-userInfoFixed .des").Text(), " ", "", -1)
	arr := strings.Split(strings.Replace(des, "\n", "", -1), "|")
	//lg.Printf("%+v", arr)
	p.JG = arr[0]
	p.Age, _ = strconv.Atoi(arr[1][:2])
	p.Education = arr[2]
	p.Marriage = arr[3]
	p.Height, _ = strconv.Atoi(arr[4][:3])
	p.Income = arr[5]
	p.Name = dom.Find(".right").Find(".nickName").Text()
	p.Gender, _ = r.Attach.(string)
	details := dom.Find(".m-userInfoDetail .m-content-box .m-btn").Nodes
	for _, v := range details {
		vl := strings.Trim(strings.Replace(v.FirstChild.Data, "\n", "", -1), " ")
		switch {
		case strings.Contains(vl, "kg"):
			if weight, err := strconv.Atoi(v.FirstChild.Data[:2]); err == nil {
				p.Weight = weight
			}
		case strings.Contains(vl, "车"):
			p.Car = vl
		case strings.Contains(vl, "房"):
			p.House = vl
		}
	}
	result.Items = append(result.Items, p)
	// find refer users
	guessLikes := dom.Find(".m-member .logo").Nodes
	for _, like := range guessLikes {
		//background-image:url(https://photo.zastatic.com/images/photo/26796/107181687/13437836295348993.png?scrop=1&crop=1&cpos=north&w=100&h=100);
		image := like.Attr[1].Val
		arr := strings.Split(image, "/")
		id := arr[len(arr)-2]
		rr := engine.Request{
			Url:       "https://album.zhenai.com/u/" + id,
			ParserFun: ParseProfile,
			Attach:    r.Attach,
		}
		result.Requests = append(result.Requests, rr)
	}
	return result
}
