package persist

import (
	"context"
	"crawler/parser/zhenai"
	"encoding/json"
	"gopkg.in/olivere/elastic.v5"
	"testing"
)

func TestItemSave(t *testing.T) {
	expected := zhenai.Profile{
		Name:      "AI",
		Gender:    "男",
		Age:       20,
		Height:    160,
		Weight:    80,
		Income:    "40000-50000",
		Education: "本科",
		Marriage:  "未婚",
		House:     "已购房",
		Car:       "已购车",
		JG:        "北京",
		Photo:     "",
	}
	id, err := save(expected)
	if err != nil {
		panic(err)
	}
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	result, err := client.Get().Index("zhenai").Type("profile").Id(id).Do(context.Background())
	if err != nil {
		panic(err)
	}
	var actual zhenai.Profile
	err = json.Unmarshal(*result.Source, &actual)
	if err != nil {
		panic(err)
	}
	if actual != expected {
		panic("'not passed")
	}

}
