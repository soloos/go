package runtime

func CurrentGLockedpoolid() int {
	_g_ := getg()
	return _g_.lockedpoolid
}

func CurrentPLockedpoolid() int {
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

func GoPoolSize(poolid int) int32 {
	return gopool.GoPoolSize(poolid)
}
