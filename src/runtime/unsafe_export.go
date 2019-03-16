package runtime

func UnsafePrintAllMStatus() {
	if sched.midle != 0 {
		println("idlem", sched.midle.ptr().id, sched.nmidle, sched.nmidlelocked)
	} else {
		println("idlem", sched.midle, sched.nmidle, sched.nmidlelocked)
	}
	println("freem", sched.freem)
	for m := sched.freem; m != nil; m = m.freelink {
		println(m.id, m.spinning, m.blocked, m.lockedInt, m.lockedg)
	}
	println("allm")
	for m := allm; m != nil; m = m.alllink {
		println(m.id, m.spinning, m.blocked, m.lockedInt, m.lockedg)
	}
}

func UnsafePrintAllPStatus() {
	for _, p := range allp {
		if p.m != 0 {
			println(p.m.ptr().id, p.id, p.lockedpoolid, p.status)
		} else {
			println("nil", p.id, p.lockedpoolid, p.status)
		}
	}
}

func UnsafeFastrand() uint32 {
	return fastrand()
}

//go:nosplit
func UnsafeProcPin() int {
	return procPin()
}

//go:nosplit
func UnsafeProcUnpin() {
	procUnpin()
}

func UnsafeGoID() int64 {
	return getg().goid
}
