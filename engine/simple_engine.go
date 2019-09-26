package engine

import (
	"crawler/fetcher"
	lg "crawler/log"
)

// request queue
var requests []Request

type SimpleEngine struct {
}

func (se *SimpleEngine) Run(seeds ...Request) {
	for _, r := range seeds {
		requests = append(requests, r)
	}
	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]
		if r.Url == "" {
			continue
		}
		work(r)
	}
}

func work(r Request) (ParseResult, error) {
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		lg.Printf("Fetcher error url %s : %v", r.Url, err)
		return ParseResult{}, err
	}
	if len(body) == 0 {
		return ParseResult{}, nil
	}
	result := r.ParserFun(&r, body)
	requests = append(requests, result.Requests...)
	/*for _, item := range result.Items {
		lg.Printf("got item %+v", item)
	}*/
	return result, nil
}
