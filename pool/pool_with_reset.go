package pool

import "sync"

type PoolWithResetFn[T any] struct {
	_pool      sync.Pool
	_resetFunc func(T)
}

func NewPoolWithResetFn[T any](newFn func() T, resetFn func(T)) *PoolWithResetFn[T] {
	return &PoolWithResetFn[T]{
		_pool: sync.Pool{
			New: func() any {
				return newFn()
			},
		},
		_resetFunc: resetFn,
	}
}

func (p *PoolWithResetFn[T]) Get() T {
	return p._pool.Get().(T)
}

func (p *PoolWithResetFn[T]) Put(v T) {
	if p._resetFunc != nil {
		p._resetFunc(v)
	}
	p._pool.Put(v)
}
