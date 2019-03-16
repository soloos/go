package runtime

// _ "runtime/cgo"

var (
	gopool gopoolt
)

const (
	maxPools        = 64
	defaultGoPoolID = 0
)

type ppoolt = []*p

type gopoolt struct {
	_excludePidsForDoRefreshDefaultPPool Bitmap

	poollock mutex

	// map poolid -> ppoolt
	ppoolReadOffset []uint32
	ppooltable      []ppoolt
}

//go:yeswritebarrierrec
func prepareGoPool() {
	gopool.ppoolReadOffset = make([]uint32, maxPools)[:1]
	gopool.ppoolReadOffset[defaultGoPoolID] = 0
	gopool.ppooltable = make([]ppoolt, maxPools)[:1]
	gopool.ppooltable[defaultGoPoolID] = gopool.ppooltable[defaultGoPoolID][:0]

	sched.poolpidle = make([]puintptr, maxPools)[:1]
	sched.poolpidle[defaultGoPoolID] = 0

	sched.poolnpidle = make([]uint32, maxPools)[:1]
	sched.poolnpidle[defaultGoPoolID] = 0

	sched.poolrunq = make([]gQueue, maxPools)[:1]
	sched.poolrunq[defaultGoPoolID] = gQueue{}

	sched.poolrunqsize = make([]int32, maxPools)[:1]
	sched.poolrunqsize[defaultGoPoolID] = 0
}

func prepareGoPPoolTable() {
	for i, _ := range allp {
		gopool.ppooltable[defaultGoPoolID] = append(gopool.ppooltable[defaultGoPoolID], allp[i])
	}
}

// call doAllocGoPool should lock(&this.poollock)
func (this *gopoolt) doAllocGoPool() int {
	var poolid = len(this.ppooltable)

	this.ppoolReadOffset = append(this.ppoolReadOffset, 0)
	this.ppoolReadOffset[poolid] = 0
	this.ppooltable = append(this.ppooltable, ppoolt{})
	this.ppooltable[poolid] = this.ppooltable[poolid][:0]

	sched.poolpidle = append(sched.poolpidle, 0)
	sched.poolpidle[poolid] = 0

	sched.poolnpidle = append(sched.poolnpidle, 0)
	sched.poolnpidle[poolid] = 0

	sched.poolrunq = append(sched.poolrunq, gQueue{})
	sched.poolrunq[poolid] = gQueue{}

	sched.poolrunqsize = append(sched.poolrunqsize, 0)
	sched.poolrunqsize[poolid] = 0

	return poolid
}

func (this *gopoolt) RefreshDefaultPool() {
	lock(&this.poollock)
	this.doRefreshDefaultPPool()
	unlock(&this.poollock)
}

func (this *gopoolt) AllocGoPool() int {
	var poolid int
	lock(&this.poollock)
	poolid = this.doAllocGoPool()
	unlock(&this.poollock)
	return poolid
}

func (this *gopoolt) RunOnGoPool(poolid int) {
	_g_ := getg().m.curg
	_g_.lockedpoolid = poolid
}

func (this *gopoolt) GoPoolRunqSize(poolid int) int32 {
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
