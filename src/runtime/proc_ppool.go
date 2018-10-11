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
	for i := 0; i < len(this._excludePidsForDoRefreshDefaultPPool); i++ {
		this._excludePidsForDoRefreshDefaultPPool[i] = false
		this._newDefaultPoolPForDoRefreshDefaultPPool[i] = nil
	}

	for _, ppool := range this.ppooltable {
		for _, p := range ppool {
			this._excludePidsForDoRefreshDefaultPPool[p.id] = true
		}
	}

	for _, p := range allp {
		if this._excludePidsForDoRefreshDefaultPPool[p.id] {
			continue
		}
		this._newDefaultPoolPForDoRefreshDefaultPPool[p.id] = p
	}

	for _, p := range this._newDefaultPoolPForDoRefreshDefaultPPool {
		if p != nil {
			this.ppooltable[defaultGoPoolID] = append(this.ppooltable[defaultGoPoolID], p)
		}
	}
}

//go:yeswritebarrierrec
func (this *gopoolt) GoPoolSetP(poolid int, pid int32) {

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
		if this.GetSchedPoolNPidle(_p_.lockedpoolid) != 0 && atomic.Load(&sched.nmspinning) == 0 {
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
