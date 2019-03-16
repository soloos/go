package main

import (
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	w     sync.WaitGroup
	count       = 0
	gs    int   = 1000000
	gDone int32 = 0

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
	atomic.AddInt32(&gDone, 1)
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
	log.Println("hello size",
		runtime.GoPoolSchednmspinning(),
		runtime.GoPoolRunqSize(-1),
		runtime.GoPoolRunqSize(gopoolDefault)+
			runtime.GoPoolRunqSize(gopool0)+
			runtime.GoPoolRunqSize(gopool1),
		runtime.GoPoolRunqSize(gopoolDefault),
		runtime.GoPoolRunqSize(gopool0),
		runtime.GoPoolRunqSize(gopool1),
	)

	runtime.UnsafePrintAllMStatus()
	runtime.UnsafePrintAllPStatus()
	log.Println("print pool done")
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	runtime.GOMAXPROCS(6)

	runtime.GoPoolSetP(gopool0, 1)
	runtime.GoPoolSetP(gopool0, 2)
	runtime.GoPoolSetP(gopool0, 3)
	runtime.GoPoolSetP(gopool1, 4)
	runtime.GoPoolSetP(gopool1, 5)

	runtime.RunOnGoPool(gopoolDefault)

	log.Println(runtime.CurrentPid(), runtime.CurrentPLockedPoolID())
	go func() {
		for {
			log.Println("start print pool", runtime.CurrentPid(), runtime.CurrentPLockedPoolID())
			printpool()
			time.Sleep(time.Second)
			log.Println("sleep done", atomic.LoadInt32(&gDone))
		}
	}()

	test1()
	test2()
	w.Wait()

	time.Sleep(time.Second * 120)
}
