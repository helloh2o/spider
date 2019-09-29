package persist

import (
	"context"
	lg "crawler/log"
	"gopkg.in/olivere/elastic.v5"
)

func ItemSaver() chan interface{} {
	in := make(chan interface{})
	go func() {
		index := 0
		for {
			item := <-in
			if _, ok := item.(string); ok {
				continue
			}
			index++
			lg.Printf("Got item #%d %v", index, item)
			if _, err := save(item); err != nil {
				lg.Printf("save item error %v", err)
			}

		}
	}()
	return in
}

func save(item interface{}) (string, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return "", err
	}
	resp, err := client.Index().Index("zhenai").Type("profile").BodyJson(item).Do(context.Background())
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}
