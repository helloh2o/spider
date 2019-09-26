package main

import (
	"crawler/engine"
	"crawler/parser/zhenai"
	"crawler/persist"
	"crawler/scheduler"
)

func main() {
	// https://h7seebfh.mirror.aliyuncs.com
	/*engine.SimpleEngine{}.Run(engine.Request{
		Url:       "https://www.zhenai.com/zhenghun",
		ParserFun: zhenai.ParseCityList,
	})*/
	seed := "https://www.zhenai.com/zhenghun"
	seedFunc := zhenai.ParseCityList

	/*seed = "http://www.zhenai.com/zhenghun/beijing"
	seedFunc = zhenai.ParseCity*/
	concurrentEngin := engine.ConcurrentEngine{
		Scheduler: &scheduler.QueuedScheduler{},
		MaxWorker: 10,
		ItemChan:  persist.ItemSaver(),
	}
	concurrentEngin.Run(engine.Request{
		Url:       seed,
		ParserFun: seedFunc,
	})
}
