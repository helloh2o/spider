package persist

import lg "crawler/log"

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
		}
	}()
	return in
}
