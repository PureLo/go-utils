package pool

import (
	"testing"
)

func BenchmarkPoolWithResetFn(b *testing.B) {
	newFn := func() []byte {
		return make([]byte, 1024)
	}
	resetFn := func(buf []byte) {
		for i := range buf {
			buf[i] = 0
		}
	}

	pool := NewPoolWithResetFn(newFn, resetFn)
	for b.Loop() {
		b := pool.Get()
		pool.Put(b)
	}
}

func TestPoolWithResetFn(t *testing.T) {
	newFn := func() []byte {
		return make([]byte, 0, 1024)
	}
	resetFn := func(buf []byte) {
		for i := range buf {
			buf[i] = 0
		}
	}
	pool := NewPoolWithResetFn(newFn, resetFn)

	for range 10 {
		b := pool.Get()
		pool.Put(b)
	}
}
