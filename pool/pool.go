package pool

import "sync"

// A Pool is a type-safe wrapper around a sync.Pool.
// Support generic
type Pool[T any] struct {
	_pool sync.Pool
}

func New[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		_pool: sync.Pool{
			New: func() any {
				return fn()
			},
		},
	}
}
func (p *Pool[T]) Get() T {
	return p._pool.Get().(T)
}

func (p *Pool[T]) Put(v T) {
	p._pool.Put(v)
}
