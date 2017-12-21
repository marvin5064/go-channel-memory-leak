package main

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type Data struct {
	ID    int
	Value string
}

func main() {

	var (
		list        []Data
		processTime = 20 * time.Minute
	)

	for i := 0; i < 50; i++ {
		item := Data{i, fmt.Sprintf("test$%v", i)}
		list = append(list, item)
	}

	m := &runtime.MemStats{}
	now := time.Now()

	fmt.Println("This application will process for", processTime)
	for {
		wg := new(sync.WaitGroup)
		Write(list, wg, true)
		wg.Wait()

		time.Sleep(200 * time.Millisecond)

		runtime.ReadMemStats(m)
		timePassed := time.Since(now)
		fmt.Printf("after %v; Alloc %d; TotalAlloc %d; Sys %d; NumGC %d; HeapAlloc %d; HeapSys %d; HeapObjects %d; HeapReleased %d;\n", timePassed, m.Alloc, m.TotalAlloc, m.Sys, m.NumGC, m.HeapAlloc, m.HeapSys, m.HeapObjects, m.HeapReleased)
		runtime.GC()
		debug.FreeOSMemory()
		if timePassed > processTime {
			break
		}
	}
	fmt.Println("Done.")
}

func Write(list []Data, wg *sync.WaitGroup, putError bool) <-chan error {
	errCh := make(chan error, 1)

	for i := range list {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			time.Sleep(200 * time.Millisecond)
			if putError {
				select {
				case errCh <- errors.New("error"):
				default:
				}
			}
		}(i, wg)
	}
	return errCh
}
