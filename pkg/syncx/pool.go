package syncx

import (
	"context"
	"sync"
)

type Pool[T any] struct {
	pool *sync.Pool
}

func NewPool[T any](new func() T) *Pool[T] {
	return &Pool[T]{
		pool: &sync.Pool{
			New: func() any {
				return new()
			},
		},
	}
}

func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}

type fabric[T any] func(ctx context.Context) (T, error)

type ErrorPool[T any] struct {
	pool   *sync.Pool
	fabric fabric[T]
}

func NewErrorPool[T any](f fabric[T]) *ErrorPool[T] {
	return &ErrorPool[T]{
		pool:   &sync.Pool{},
		fabric: f,
	}
}

func (p *ErrorPool[T]) Get(ctx context.Context) (T, error) {
	got := p.pool.Get()
	if got != nil {
		return got.(T), nil
	}

	return p.fabric(ctx)
}

func (p *ErrorPool[T]) Put(x T) {
	p.pool.Put(x)
}
