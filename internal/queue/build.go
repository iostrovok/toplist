package queue

import (
	"sync"
	"sync/atomic"
)

type Build struct {
	mx        sync.Mutex
	buildFunc func()
	threshold uint32
	CountRun  uint32

	Counter atomic.Uint32
	IsWork  atomic.Uint32
}

func NewBuild(threshold uint32, buildFunc func()) *Build {
	return &Build{
		threshold: threshold,
		buildFunc: buildFunc,
	}
}

func (b *Build) IncrementAndRun() {
	if b.Counter.Add(1) < b.threshold || b.IsWork.Load() > 0 {
		return
	}

	b.mx.Lock()
	if b.IsWork.Load() > 0 {
		b.mx.Unlock()
		return
	}

	b.IsWork.Store(1)
	b.Counter.Store(0)

	b.mx.Unlock()

	b.CountRun++
	b.buildFunc()
	b.IsWork.Store(0)
}
