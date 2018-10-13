package runtime

import "runtime/internal/atomic"

func GetAllPids() []int32 {
	var (
		result []int32
		pp     *p
	)
	for _, pp = range allp {
		result = append(result, pp.id)
	}
	return result
}

//go:yeswritebarrierrec
func GoPoolSetP(poolid int, pid int32) {
	gopool.GoPoolSetP(poolid, pid)
}

// call doRefreshDefaultPPool should lock(&gopool.poollock)
func (gopool *gopoolt) doRefreshDefaultPPool() {
	excludePids := make(map[int32]bool, 255)
	for _, ppool := range gopool.ppoolmap {
		for _, p := range *ppool {
			excludePids[p.id] = true
		}
	}

	var (
		newDefaultPoolP []*p
		exists          bool
	)
	for _, p := range allp {
		if _, exists = excludePids[p.id]; exists {
			continue
		}
		newDefaultPoolP = append(newDefaultPoolP, p)
	}

	originDefaultPool := gopool.ppoolmap[0]
	originDefaultPoolCap := cap(*originDefaultPool)
	if originDefaultPoolCap < len(newDefaultPoolP) {
		*originDefaultPool = (*originDefaultPool)[:originDefaultPoolCap]
		copy(*originDefaultPool, newDefaultPoolP)
		*originDefaultPool = append(*originDefaultPool, newDefaultPoolP[originDefaultPoolCap:]...)
	} else {
		copy(*originDefaultPool, newDefaultPoolP)
	}
}

//go:yeswritebarrierrec
func (gopool *gopoolt) GoPoolSetP(poolid int, pid int32) {
	ppool := gopool.ppoolmap[poolid]

	var pp *p
	for _, tmpp := range allp {
		if tmpp.id == pid {
			pp = tmpp
			break
		}
	}

	if pp == nil {
		panic("p not found")
	}

	if pp.lockedpoolid == poolid {
		return
	}

	lock(&gopool.poollock)
	*ppool = append(*ppool, pp)
	pp.lockedpoolid = poolid
	gopool.doRefreshDefaultPPool()
	unlock(&gopool.poollock)
}

func (gopool *gopoolt) runqput(_p_ *p, gp *g, next bool) {
	if _p_.lockedpoolid == gp.lockedpoolid {
		runqput(_p_, gp, next)
		if gopool.GetSchedPoolNPidle(_p_.lockedpoolid) != 0 && atomic.Load(&sched.nmspinning) == 0 {
			wakep(_p_.lockedpoolid)
		}
	} else {
		_p_ = gopool.chooseP(gp.lockedpoolid)
		lock(&sched.lock)
		globrunqput(gp)
		unlock(&sched.lock)
		startm(_p_.lockedpoolid, nil, false)
	}
}

// call doChooseP should lock(&gopool.poollock)
func (gopool *gopoolt) chooseP(poolid int) (result *p) {
	ppool := gopool.ppoolmap[poolid]
	result = (*ppool)[int(atomic.Xadd(gopool.ppoolReadOffset[poolid], 1)%uint32(len(*ppool)))]
	return result
}
