package pool

import "sync"

// A Pool is a type-safe wrapper around a sync.Pool.
// Support generic
type Pool[T any] struct {
	pool sync.Pool
}

func New[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return fn()
			},
		},
	}
}
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(v T) {
	p.pool.Put(v)
}
