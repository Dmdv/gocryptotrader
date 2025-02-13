package alert

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	wait := Notice{}
	var wg sync.WaitGroup

	// standard alert
	wg.Add(100)
	for x := 0; x < 100; x++ {
		go func() {
			w := wait.Wait(nil)
			wg.Done()
			if <-w {
				log.Fatal("incorrect routine wait response for alert expecting false")
			}
			wg.Done()
		}()
	}

	wg.Wait()
	wg.Add(100)
	isLeaky(&wait, nil, t)
	wait.Alert()
	wg.Wait()
	isLeaky(&wait, nil, t)

	// use kick
	ch := make(chan struct{})
	wg.Add(100)
	for x := 0; x < 100; x++ {
		go func() {
			w := wait.Wait(ch)
			wg.Done()
			if !<-w {
				log.Fatal("incorrect routine wait response for kick expecting true")
			}
			wg.Done()
		}()
	}
	wg.Wait()
	wg.Add(100)
	isLeaky(&wait, ch, t)
	close(ch)
	wg.Wait()
	ch = make(chan struct{})
	isLeaky(&wait, ch, t)

	// late receivers
	wg.Add(100)
	for x := 0; x < 100; x++ {
		go func(x int) {
			bb := wait.Wait(ch)
			wg.Done()
			if x%2 == 0 {
				time.Sleep(time.Millisecond * 5)
			}
			b := <-bb
			if b {
				log.Fatal("incorrect routine wait response since we call alert below; expecting false")
			}
			wg.Done()
		}(x)
	}
	wg.Wait()
	wg.Add(100)
	isLeaky(&wait, ch, t)
	wait.Alert()
	wg.Wait()
	isLeaky(&wait, ch, t)
}

// isLeaky tests to see if the wait functionality is returning an abnormal
// channel that is operational when it shouldn't be.
func isLeaky(a *Notice, ch chan struct{}, t *testing.T) {
	t.Helper()
	check := a.Wait(ch)
	time.Sleep(time.Millisecond * 5) // When we call wait a routine for hold is
	// spawned, so for a test we need to add in a time for goschedular to allow
	// routine to actually wait on the forAlert and kick channels
	select {
	case <-check:
		t.Fatal("leaky waiter")
	default:
	}
}
