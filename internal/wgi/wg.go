package wgi

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var mutex sync.Mutex
var maxRoutines = 9999
var routines = 0

func init() {
	routines = runtime.NumGoroutine()
}

type WaitGroupInterface interface {
	Add(delta int)
	Done()
	Wait()
}

type WaitGroup struct {
	wg sync.WaitGroup
}

func (w *WaitGroup) Add(delta int) {
	mutex.Lock()
	defer mutex.Unlock()

	// waiting for our turn
	for routines > maxRoutines {
		time.Sleep(time.Millisecond * time.Duration(10+rand.Int63n(10)))
	}

	w.wg.Add(delta)
	routines += delta
}

func (w *WaitGroup) Done() {
	w.wg.Done()
	routines--
}

func (w *WaitGroup) Wait() {
	w.wg.Wait()
}
