package runtime

//go:yeswritebarrierrec
func GoPoolSetP(poolid int, pid int32) {
	gopool.GoPoolSetP(poolid, pid)
}
