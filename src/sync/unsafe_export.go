package sync

import (
	"internal/race"
	"unsafe"
)

const UnsafeRaceEnabled = race.Enabled

func UnsafeRaceEnable() {
	race.Enable()
}

func UnsafeRaceDisable() {
	race.Disable()
}

func UnsafeRaceReleaseMerge(addr unsafe.Pointer) {
	race.ReleaseMerge(addr)
}

func UnsafeRaceAcquire(addr unsafe.Pointer) {
	race.Acquire(addr)
}
