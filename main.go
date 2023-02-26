package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var a int = 0
	var wg sync.WaitGroup
	var ch = make(chan struct{})
	go func() {
		wg.Add(1)
		defer wg.Done()
		for a < 10 {
			select {
			case <-ch:
				fmt.Println("received")
				return
			default:
				a++
				time.Sleep(time.Second)
			}
		}
	}()
	_ = WaitForNum(20, &a, 5)
	//ch <- struct{}{}
	close(ch)
	wg.Wait()
}

func WaitForNum(timeout int, num *int, expected int) error {
	return MyWaitFor(timeout, func() (bool, error) {
		fmt.Println("This time, value is, ", *num)
		return *num >= expected, nil
	})
}

func MyWaitFor(timeoutSec int, predicate func() (bool, error)) error {
	var success bool
	var err error

	start := time.Now().Unix()

	for {
		// our expected state does not appear before timeout
		if timeoutSec > 0 && time.Now().Unix()-start > int64(timeoutSec) {
			return fmt.Errorf("MyWaitFor timeout")
		}
		time.Sleep(1 * time.Second)

		ch := make(chan struct{}, 0)
		go func() {
			defer close(ch)
			success, err = predicate()
		}()

		select {
		case <-ch:
			if err != nil {
				return fmt.Errorf("error in MyWaitFor [%w]", err)
			}
			if success {
				return nil
			}
		case <-time.After(time.Duration(timeoutSec) * time.Second):
			return fmt.Errorf("MyWaitFor a predicate does not return before timeout")
		}
	}
}
