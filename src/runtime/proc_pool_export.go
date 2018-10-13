package runtime

func CurrentGLockedpoolid() int {
	_g_ := getg()
	return _g_.lockedpoolid
}

func CurrentPLockedpoolid() int {
	_g_ := getg()
	_p_ := _g_.m.p.ptr()
	return _p_.lockedpoolid
}

func CurrentPid() int32 {
	_g_ := getg()
	_p_ := _g_.m.p.ptr()
	return _p_.id
}

func AllocGoPool(poolid int) {
	gopool.AllocGoPool(poolid)
}

func SetGoPoolID(poolid int) {
	gopool.SetGoPoolID(poolid)
}

func GoPoolSize(poolid int) int32 {
	return gopool.GoPoolSize(poolid)
}
