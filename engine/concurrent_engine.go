package engine

import (
	lg "crawler/log"
)

type Scheduler interface {
	ReadyNotifier
	Submit(r Request)
	WorkerChan() chan Request
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

type ConcurrentEngine struct {
	Scheduler Scheduler
	MaxWorker int
	ItemChan  chan interface{}
}

func (engine *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult)
	// run scheduler
	engine.Scheduler.Run()
	// create workers
	for i := 0; i < engine.MaxWorker; i++ {
		engine.createWorker(engine.Scheduler.WorkerChan(), out, engine.Scheduler)
	}
	// submit req
	for _, req := range seeds {
		engine.Scheduler.Submit(req)
	}
	repeatMap := make(map[string]bool)
	// wait result from out
	for {
		result := <-out
		for _, item := range result.Items {
			go func() { engine.ItemChan <- item }()
		}
		for _, req := range result.Requests {
			if _, ok := repeatMap[req.Url]; ok {
				//lg.Printf("Reap URL %s", req.Url)
				continue
			} else {
				engine.Scheduler.Submit(req)
				repeatMap[req.Url] = true
			}
		}
	}
}

func (engine *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, notifier ReadyNotifier) {
	go func() {
		// loop
		for {
			// worker ready
			notifier.WorkerReady(in)
			r := <-in
			result, err := work(r)
			if err != nil {
				lg.Printf("work error %v", err)
			} else {
				out <- result
			}
		}
	}()
}
