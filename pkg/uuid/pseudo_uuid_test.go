package uuid

import (
	"sync"
	"testing"
)

func Benchmark_random(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		s := uint32(0)
		for pb.Next() {
			s += random()
		}
		_ = s // to prevent optimizer from removing the loop
	})
}

func Benchmark_randomNoPool(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		pcg := newPCG()
		s := uint32(0)
		for p.Next() {
			s += pcg.random()
		}
		_ = s // to prevent optimizer from removing the loop
	})
}

func Benchmark_randomWithLock(b *testing.B) {
	pcg := newPCG()
	mu := sync.Mutex{}
	b.RunParallel(func(p *testing.PB) {
		s := uint32(0)
		for p.Next() {
			mu.Lock()
			s += pcg.random()
			mu.Unlock()
		}
		_ = s // to prevent optimizer from removing the loop
	})
}
