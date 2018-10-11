package runtime

import "runtime/internal/atomic"

const (
	maxPools      = 255
	maxProcessors = 2048
)

const (
	defaultGoPoolID = 0
)

type ppoolt = []*p

type gopoolt struct {
	_excludePidsForDoRefreshDefaultPPool     [maxProcessors]bool
	_newDefaultPoolPForDoRefreshDefaultPPool [maxProcessors]*p

	glistlens []int

	poollock mutex

	// map poolid -> ppoolt
	ppoolReadOffset  []uint32
	_ppooltable_pool [maxPools]ppoolt
	ppooltable       []ppoolt
}

//go:yeswritebarrierrec
func prepareGoPool() {
	gopool.glistlens = make([]int, maxPools)[:1]
	gopool.glistlens[defaultGoPoolID] = 0

	gopool.ppoolReadOffset = make([]uint32, maxPools)[:1]
	gopool.ppoolReadOffset[defaultGoPoolID] = 0
	gopool.ppooltable = make([]ppoolt, maxPools)[:1]
	gopool.ppooltable[defaultGoPoolID] = gopool.ppooltable[defaultGoPoolID][:0]

	sched.poolpidle = make([]puintptr, maxPools)[:1]
	sched.poolpidle[defaultGoPoolID] = 0

	sched.poolnpidle = make([]uint32, maxPools)[:1]
	sched.poolnpidle[defaultGoPoolID] = 0

	sched.poolrunqhead = make([]guintptr, maxPools)[:1]
	sched.poolrunqhead[defaultGoPoolID] = 0

	sched.poolrunqtail = make([]guintptr, maxPools)[:1]
	sched.poolrunqtail[defaultGoPoolID] = 0

	sched.poolrunqsize = make([]int32, maxPools)[:1]
	sched.poolrunqsize[defaultGoPoolID] = 0
}

func prepareGoDefaultPool() {
	for i, _ := range allp {
		gopool.ppooltable[defaultGoPoolID] = append(gopool.ppooltable[defaultGoPoolID], allp[i])
	}
}

// call doRefreshDefaultPool should lock(&this.poollock)
func (this *gopoolt) doRefreshDefaultPool() {
	this.doRefreshDefaultPPool()
}

// call doAllocGoPool should lock(&this.poollock)
func (this *gopoolt) doAllocGoPool() int {
	var poolid = len(this.ppooltable)

	this.glistlens = append(this.glistlens, 0)
	this.glistlens[poolid] = 0

	this.ppoolReadOffset = append(this.ppoolReadOffset, 0)
	this.ppoolReadOffset[poolid] = 0
	this.ppooltable = append(this.ppooltable, ppoolt{})
	this.ppooltable[poolid] = this.ppooltable[poolid][:0]

	sched.poolpidle = append(sched.poolpidle, 0)
	sched.poolpidle[poolid] = 0

	sched.poolnpidle = append(sched.poolnpidle, 0)
	sched.poolnpidle[poolid] = 0

	sched.poolrunqhead = append(sched.poolrunqhead, 0)
	sched.poolrunqhead[poolid] = 0

	sched.poolrunqtail = append(sched.poolrunqtail, 0)
	sched.poolrunqtail[poolid] = 0

	sched.poolrunqsize = append(sched.poolrunqsize, 0)
	sched.poolrunqsize[poolid] = 0

	return poolid
}

func (this *gopoolt) RefreshDefaultPool() {
	lock(&this.poollock)
	this.doRefreshDefaultPool()
	unlock(&this.poollock)
}

func (this *gopoolt) AllocGoPool() int {
	var poolid int
	lock(&this.poollock)
	poolid = this.doAllocGoPool()
	unlock(&this.poollock)
	return poolid
}

//go:nowritebarrierrec
func (this *gopoolt) PushSchedPoolPidle(poolid int, _p_ *p) {
	sched.poolpidle[poolid].set(_p_)
	atomic.Xadd(&sched.poolnpidle[poolid], 1)
	atomic.Xadd(&sched.poolsnpidle, 1) // TODO: fast atomic
}

//go:nowritebarrierrec
func (this *gopoolt) PopSchedPoolPidle(poolid int) *p {
	_p_ := sched.poolpidle[poolid].ptr()
	if _p_ != nil {
		sched.poolpidle[poolid].set(_p_.link.ptr())
		atomic.Xadd(&sched.poolnpidle[poolid], -1)
		atomic.Xadd(&sched.poolsnpidle, -1) // TODO: fast atomic
	}
	return _p_
}

//go:nowritebarrierrec
func (this *gopoolt) GetSchedPoolPidle(poolid int) *puintptr {
	return &sched.poolpidle[poolid]
}

//go:nowritebarrierrec
func (this *gopoolt) GetSchedPoolNPidle(poolid int) uint32 {
	return atomic.Load(&sched.poolnpidle[poolid])
}

//go:nowritebarrierrec
func (this *gopoolt) GetSchedPoolsNPidle() uint32 {
	return atomic.Load(&sched.poolsnpidle)
}

//go:nowritebarrierrec
func (this *gopoolt) GetSchedRunqhead(poolid int) *guintptr {
	return &sched.poolrunqhead[poolid]
}

//go:nowritebarrierrec
func (this *gopoolt) GetSchedRunqtail(poolid int) *guintptr {
	return &sched.poolrunqtail[poolid]
}

//go:nowritebarrierrec
func (this *gopoolt) GetSchedRunqsize(poolid int) *int32 {
	return &sched.poolrunqsize[poolid]
}

//go:nowritebarrierrec
func (this *gopoolt) RunOnGoPool(poolid int) {
	_g_ := getg()
	_g_.lockedpoolid = poolid
}

//go:nowritebarrierrec
func (this *gopoolt) GoPoolSize(poolid int) int32 {
	lock(&this.poollock)
	var result int32 = 0
	if poolid == -1 {
		for _, size := range sched.poolrunqsize {
			result += size
		}
	} else {
		result = sched.poolrunqsize[poolid]
	}
	unlock(&this.poollock)
	return result
}
