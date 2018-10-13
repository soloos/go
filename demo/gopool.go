package main

import (
	"log"
	"runtime"
	"sync"
	"time"
)

var (
	w     sync.WaitGroup
	count = 0
	gs    = 1000000
)

func basetest() {
	max := 10000
	for i := 0; i < max; i++ {
		count += i
	}
	w.Done()
}

func test1() {
	go func() {
		runtime.SetGoPoolID(1)
		for i := 0; i < gs; i++ {
			w.Add(1)
			go basetest()
		}
	}()
}

func test2() {
	go func() {
		runtime.SetGoPoolID(2)
		for i := 0; i < gs; i++ {
			w.Add(1)
			go basetest()
		}
	}()
}

func printpool() {
	if runtime.GoPoolSize(-1) > 0 {
		log.Println("hello -1", runtime.GoPoolSize(-1))
	}
	if runtime.GoPoolSize(0) > 0 {
		log.Println("hello 0", runtime.GoPoolSize(0))
	}
	if runtime.GoPoolSize(1) > 0 {
		log.Println("hello 1", runtime.GoPoolSize(1))
	}
	if runtime.GoPoolSize(2) > 0 {
		log.Println("hello 2", runtime.GoPoolSize(2))
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	runtime.GOMAXPROCS(6)

	runtime.AllocGoPool(1)
	runtime.AllocGoPool(2)

	runtime.GoPoolSetP(1, 1)
	runtime.GoPoolSetP(1, 2)
	runtime.GoPoolSetP(2, 3)
	runtime.GoPoolSetP(2, 4)
	runtime.GoPoolSetP(2, 5)

	runtime.SetGoPoolID(0)

	log.Println(runtime.CurrentPid(), runtime.CurrentPLockedpoolid())
	go func() {
		for {
			log.Println("start print pool", runtime.CurrentPid(), runtime.CurrentPLockedpoolid())
			printpool()
			time.Sleep(time.Second)
			log.Println("sleep done")
		}
	}()

	test1()
	test2()
	w.Wait()
	log.Println(count)

	time.Sleep(time.Second * 100)
}
