package spinlock

import (
	"runtime"
	"sync/atomic"
)

type Lock struct {
	lock uintptr
}

func (l *Lock) Lock() {
	for !atomic.CompareAndSwapUintptr(&l.lock, 0, 1) {
		runtime.Gosched()
	}
}

func (l *Lock) Unlock() {
	atomic.StoreUintptr(&l.lock, 0)
}
