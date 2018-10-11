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

	gopoolDefault = 0
	gopool0       = runtime.AllocGoPool()
	gopool1       = runtime.AllocGoPool()
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
		runtime.RunOnGoPool(gopool0)
		for i := 0; i < gs; i++ {
			w.Add(1)
			go basetest()
		}
	}()
}

func test2() {
	go func() {
		runtime.RunOnGoPool(gopool1)
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
	if runtime.GoPoolSize(gopoolDefault) > 0 {
		log.Println("hello 0", runtime.GoPoolSize(gopoolDefault))
	}
	if runtime.GoPoolSize(gopool0) > 0 {
		log.Println("hello 1", runtime.GoPoolSize(gopool0))
	}
	if runtime.GoPoolSize(gopool1) > 0 {
		log.Println("hello 2", runtime.GoPoolSize(gopool1))
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	runtime.GOMAXPROCS(6)

	runtime.GoPoolSetP(gopool0, 1)
	runtime.GoPoolSetP(gopool0, 2)
	runtime.GoPoolSetP(gopool1, 3)
	runtime.GoPoolSetP(gopool1, 4)
	runtime.GoPoolSetP(gopool1, 5)

	runtime.RunOnGoPool(gopoolDefault)

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
