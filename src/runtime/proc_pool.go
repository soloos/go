package runtime

import "runtime/internal/atomic"

type ppoolt = []*p

type gopoolt struct {
	poollock mutex

	// map poolid -> ppoolt
	ppoolReadOffset map[int]*uint32
	ppoolmap        map[int]*ppoolt
}

func prepareGoPool() {
	gopool.ppoolmap = make(map[int]*ppoolt, 255)
	gopool.ppoolReadOffset = make(map[int]*uint32, 255)

	sched.poolpidle = make([]*puintptr, 255)
	sched.poolnpidle = make([]*uint32, 255)

	sched.poolrunqhead = make(map[int]*guintptr, 255)
	sched.poolrunqtail = make(map[int]*guintptr, 255)
	sched.poolrunqsize = make(map[int]*int32, 255)

	gopool.doAllocGoPool(0)
	gopool.doRefreshDefaultPool()
}

// call doRefreshDefaultPool should lock(&gopool.poollock)
func (gopool *gopoolt) doRefreshDefaultPool() {
	gopool.doRefreshDefaultPPool()
}

// call doAllocGoPool should lock(&gopool.poollock)
func (gopool *gopoolt) doAllocGoPool(poolid int) {
	gopool.ppoolmap[poolid] = new(ppoolt)
	gopool.ppoolReadOffset[poolid] = new(uint32)

	sched.poolpidle[poolid] = new(puintptr)
	*sched.poolpidle[poolid] = 0

	sched.poolnpidle[poolid] = new(uint32)
	*sched.poolnpidle[poolid] = 0

	sched.poolrunqhead[poolid] = new(guintptr)
	*sched.poolrunqhead[poolid] = 0

	sched.poolrunqtail[poolid] = new(guintptr)
	*sched.poolrunqtail[poolid] = 0

	sched.poolrunqsize[poolid] = new(int32)
}

func (gopool *gopoolt) RefreshDefaultPool() {
	lock(&gopool.poollock)
	gopool.doRefreshDefaultPool()
	unlock(&gopool.poollock)
}

func (gopool *gopoolt) AllocGoPool(poolid int) {
	lock(&gopool.poollock)
	gopool.doAllocGoPool(poolid)
	unlock(&gopool.poollock)
}

//go:nowritebarrierrec
func (gopool *gopoolt) PushSchedPoolPidle(poolid int, _p_ *p) {
	sched.poolpidle[poolid].set(_p_)
	atomic.Xadd(sched.poolnpidle[poolid], 1)
	atomic.Xadd(&sched.poolsnpidle, 1) // TODO: fast atomic
}

//go:nowritebarrierrec
func (gopool *gopoolt) PopSchedPoolPidle(poolid int) *p {
	_p_ := sched.poolpidle[poolid].ptr()
	if _p_ != nil {
		sched.poolpidle[poolid].set(_p_.link.ptr())
		atomic.Xadd(sched.poolnpidle[poolid], -1)
		atomic.Xadd(&sched.poolsnpidle, -1) // TODO: fast atomic
	}
	return _p_
}

//go:nowritebarrierrec
func (gopool *gopoolt) GetSchedPoolPidle(poolid int) *puintptr {
	return sched.poolpidle[poolid]
}

//go:nowritebarrierrec
func (gopool *gopoolt) GetSchedPoolNPidle(poolid int) uint32 {
	return atomic.Load(sched.poolnpidle[poolid])
}

//go:nowritebarrierrec
func (gopool *gopoolt) GetSchedPoolsNPidle() uint32 {
	return atomic.Load(&sched.poolsnpidle)
}

//go:nowritebarrierrec
func (gopool *gopoolt) GetSchedRunqhead(poolid int) *guintptr {
	return sched.poolrunqhead[poolid]
}

//go:nowritebarrierrec
func (gopool *gopoolt) GetSchedRunqtail(poolid int) *guintptr {
	return sched.poolrunqtail[poolid]
}

//go:nowritebarrierrec
func (gopool *gopoolt) GetSchedRunqsize(poolid int) *int32 {
	return sched.poolrunqsize[poolid]
}

//go:nowritebarrierrec
func (gopool *gopoolt) SetGoPoolID(poolid int) {
	_g_ := getg()
	_g_.lockedpoolid = poolid
}

//go:nowritebarrierrec
func (gopool *gopoolt) GoPoolSize(poolid int) int32 {
	lock(&gopool.poollock)
	var result int32 = 0
	if poolid == -1 {
		for _, size := range sched.poolrunqsize {
			result += *size
		}
	} else {
		result = *sched.poolrunqsize[poolid]
	}
	unlock(&gopool.poollock)
	return result
}
