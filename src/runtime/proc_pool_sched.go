package runtime

import "runtime/internal/atomic"

func (this *gopoolt) GetSchednmspinning() uint32 {
	return atomic.Load(&sched.nmspinning)
}

func (this *gopoolt) GetSchedPoolPidle(poolid int) *puintptr {
	return &sched.poolpidle[poolid]
}

func (this *gopoolt) GetSchedPoolNPidle(poolid int) uint32 {
	return atomic.Load(&sched.poolnpidle[poolid])
}

func (this *gopoolt) GetSchedPoolNPidleSum() uint32 {
	return atomic.Load(&sched.poolnpidle_sum)
}

func (this *gopoolt) GetSchedRunq(poolid int) *gQueue {
	return &sched.poolrunq[poolid]
}

func (this *gopoolt) GetSchedRunqSize(poolid int) *int32 {
	return &sched.poolrunqsize[poolid]
}

func (this *gopoolt) PushSchedPoolPidle(poolid int, _p_ *p) {
	sched.poolpidle[poolid].set(_p_)
	atomic.Xadd(&sched.poolnpidle[poolid], 1)
	atomic.Xadd(&sched.poolnpidle_sum, 1) // TODO: fast atomic
}

func (this *gopoolt) PopSchedPoolPidle(poolid int) *p {
	_p_ := sched.poolpidle[poolid].ptr()
	if _p_ != nil {
		sched.poolpidle[poolid].set(_p_.link.ptr())
		atomic.Xadd(&sched.poolnpidle[poolid], -1)
		atomic.Xadd(&sched.poolnpidle_sum, -1) // TODO: fast atomic
	}
	return _p_
}

func (this *gopoolt) StartAllSchedPoolPidle(n int32, spinning bool) {
	for poolid, _ := range gopool.ppooltable {
		for ; n != 0 && gopool.GetSchedPoolNPidle(poolid) != 0; n-- {
			startm(poolid, nil, spinning)
		}
	}
}
