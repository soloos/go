package runtime

import "runtime/internal/atomic"

// P (processor) pool

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

// call doRefreshDefaultPPool should lock(&this.poollock)
func (this *gopoolt) doRefreshDefaultPPool() {
	this._excludePidsForDoRefreshDefaultPPool.Reset()

	for i := 1; i < len(this.ppooltable); i++ {
		ppool := this.ppooltable[i]
		for _, p := range ppool {
			this._excludePidsForDoRefreshDefaultPPool.Set(p.id)
		}
	}

	this.ppooltable[defaultGoPoolID] = this.ppooltable[defaultGoPoolID][:0]
	for _, p := range allp {
		if this._excludePidsForDoRefreshDefaultPPool.Has(p.id) {
			continue
		}
		this.ppooltable[defaultGoPoolID] = append(this.ppooltable[defaultGoPoolID], p)
	}
}

func (this *gopoolt) GoPPoolSetP(poolid int, pid int32) {

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

	lock(&this.poollock)
	this.ppooltable[poolid] = append(this.ppooltable[poolid], pp)
	pp.lockedpoolid = poolid
	this.doRefreshDefaultPPool()
	unlock(&this.poollock)
}

func (this *gopoolt) runqput(_p_ *p, gp *g, next bool) {
	if _p_.lockedpoolid == gp.lockedpoolid {
		runqput(_p_, gp, next)
		if this.GetSchedPoolNPidle(_p_.lockedpoolid) != 0 && gopool.GetSchednmspinning(_p_.lockedpoolid) == 0 {
			wakep(_p_.lockedpoolid)
		}
	} else {
		_p_ = this.chooseP(gp.lockedpoolid)
		lock(&sched.lock)
		globrunqput(gp)
		unlock(&sched.lock)
		startm(_p_.lockedpoolid, nil, false)
	}
}

// call doChooseP should lock(&this.poollock)
func (this *gopoolt) chooseP(poolid int) (result *p) {
	ppool := this.ppooltable[poolid]
	choosenID := int(atomic.Xadd(&this.ppoolReadOffset[poolid], 1) % uint32(len(ppool)))
	return ppool[choosenID]
}
