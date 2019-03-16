package runtime

func CurrentGLockedPoolID() int {
	_g_ := getg()
	return _g_.lockedpoolid
}

func CurrentPLockedPoolID() int {
	_p_ := getg().m.p.ptr()
	return _p_.lockedpoolid
}

func CurrentPid() int32 {
	_p_ := getg().m.p.ptr()
	return _p_.id
}

func AllocGoPool() int {
	return gopool.AllocGoPool()
}

func RunOnGoPool(poolid int) {
	gopool.RunOnGoPool(poolid)
}

func GoPoolSchednmspinning() uint32 {
	return gopool.GetSchednmspinning()
}

func GoPoolRunqSize(poolid int) int32 {
	return gopool.GoPoolRunqSize(poolid)
}

func GoPoolSetP(poolid int, pid int32) {
	gopool.GoPPoolSetP(poolid, pid)
}
