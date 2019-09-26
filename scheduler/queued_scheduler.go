package scheduler

import "crawler/engine"

type QueuedScheduler struct {
	requestChan chan engine.Request
	workerChan  chan chan engine.Request
}

func (q *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

func (q *QueuedScheduler) Submit(r engine.Request) {
	q.requestChan <- r
}
func (q *QueuedScheduler) WorkerReady(w chan engine.Request) {
	q.workerChan <- w
}
func (q *QueuedScheduler) Run() {
	q.requestChan = make(chan engine.Request)
	q.workerChan = make(chan chan engine.Request)
	go func() {
		var requestQ []engine.Request
		var workerQ []chan engine.Request
		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeRequest = requestQ[0]
				activeWorker = workerQ[0]
			}
			select {
			case request := <-q.requestChan:
				requestQ = append(requestQ, request)
			case worker := <-q.workerChan:
				workerQ = append(workerQ, worker)
			// send active request to active worker
			case activeWorker <- activeRequest:
				// update queue
				requestQ = requestQ[1:]
				workerQ = workerQ[1:]
			}
		}
	}()
}
func (q QueuedScheduler) ConfigureReqChan(chan engine.Request) {
}
